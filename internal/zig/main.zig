const std = @import("std");
const dynamics = @import("dynamics.zig");

// Internal PRNG
const RndGen = struct {
    state: u32,
    pub fn init(seed: u32) RndGen {
        return RndGen{ .state = if (seed == 0) 0xDEADBEEF else seed };
    }
    pub fn random(self: *RndGen) Random {
        return Random{ .ptr = self };
    }
    pub fn next(self: *RndGen) u32 {
        var x = self.state;
        x ^= x << 13;
        x ^= x >> 17;
        x ^= x << 5;
        self.state = x;
        return x;
    }
    const Random = struct {
        ptr: *RndGen,
        pub fn intRangeAtMost(self: Random, T: type, max: T) T {
            const r = self.ptr.next();
            return @intCast(r % (max + 1));
        }
        pub fn intRangeLessThan(self: Random, T: type, max: T) T {
             const r = self.ptr.next();
             return @intCast(r % max);
        }
        pub fn float(self: Random, T: type) T {
            const r = self.ptr.next();
            return @as(T, @floatFromInt(r)) / @as(T, @floatFromInt(std.math.maxInt(u32)));
        }
    };
};

const ALPHABET_SIZE = 2;
const MAX_ENERGY = 100.0;
const MIN_ENERGY = 0.1;
const INITIAL_ENERGY = 1.0;

// Config Struct (C-ABI Compatible)
pub const Config = extern struct {
    pop_size: i32,
    path_len: i32,
    alpha: f32, // Learning Rate
    beta: f32,  // Decay
    gamma: f32, // Turnover Threshold
    seed: u64,
};

const Path = struct {
    id: u32,
    energy: f32,
    symbols: []u8, // Dynamic length now supported via slice

    pub fn init(allocator: std.mem.Allocator, id: u32, length: usize, rnd: *RndGen.Random) !Path {
        var p = Path{
            .id = id,
            .energy = INITIAL_ENERGY,
            .symbols = try allocator.alloc(u8, length),
        };
        for (p.symbols) |*sym| {
            sym.* = @intCast(rnd.intRangeAtMost(u8, ALPHABET_SIZE - 1));
        }
        return p;
    }

    pub fn deinit(self: *Path, allocator: std.mem.Allocator) void {
        allocator.free(self.symbols);
    }

    pub fn mutate(self: *Path, rnd: *RndGen.Random) void {
        if (self.symbols.len == 0) return;
        const idx = rnd.intRangeLessThan(usize, self.symbols.len);
        self.symbols[idx] = @intCast(rnd.intRangeAtMost(u8, ALPHABET_SIZE - 1));
        self.energy = INITIAL_ENERGY;
    }
};

const Population = struct {
    paths: []Path,
    allocator: std.mem.Allocator,
    prng: RndGen,
    config: Config,

    pub fn init(allocator: std.mem.Allocator, cfg: Config) !*Population {
        var pop = try allocator.create(Population);
        pop.allocator = allocator;
        pop.config = cfg;
        pop.prng = RndGen.init(@truncate(cfg.seed));
        pop.paths = try allocator.alloc(Path, @intCast(cfg.pop_size));
        
        var rnd = pop.prng.random();
        for (pop.paths, 0..) |*p, i| {
            p.* = try Path.init(allocator, @intCast(i), @intCast(cfg.path_len), &rnd);
        }
        return pop;
    }

    pub fn deinit(self: *Population) void {
        for (self.paths) |*p| {
            p.deinit(self.allocator);
        }
        self.allocator.free(self.paths);
        self.allocator.destroy(self);
    }

    pub fn select(self: *Population) u32 {
        var total_energy: f32 = 0.0;
        for (self.paths) |p| {
            total_energy += p.energy;
        }

        var rnd = self.prng.random();
        const threshold = rnd.float(f32) * total_energy;
        
        var current: f32 = 0.0;
        for (self.paths) |p| {
            current += p.energy;
            if (current >= threshold) {
                return p.id;
            }
        }
        return self.paths[0].id;
    }

    pub fn update(self: *Population, selected_id: u32, reward: f32) void {
        // Mean
        var total_energy: f32 = 0.0;
        for (self.paths) |p| {
            total_energy += p.energy;
        }
        const mean_energy = total_energy / @as(f32, @floatFromInt(self.paths.len));

        // Update
        for (self.paths) |*p| {
            if (p.id == selected_id) {
                // Modular Dynamics Call
                p.energy = dynamics.update_winner(p.energy, self.config.alpha, reward, mean_energy);
            } else {
                // Modular Dynamics Call
                p.energy = dynamics.decay_loser(p.energy, self.config.beta);
            }

            if (p.energy > MAX_ENERGY) p.energy = MAX_ENERGY;
            if (p.energy < MIN_ENERGY) p.energy = MIN_ENERGY;
        }

        // Turnover
        var rnd = self.prng.random();
        var best_path_idx: usize = 0;
        var max_e: f32 = -1.0;
        for (self.paths, 0..) |p, i| {
            if (p.energy > max_e) {
                max_e = p.energy;
                best_path_idx = i;
            }
        }

        for (self.paths) |*p| {
            if (p.energy <= MIN_ENERGY + self.config.gamma) { // Configurable threshold
                @memcpy(p.symbols, self.paths[best_path_idx].symbols);
                p.mutate(&rnd);
            }
        }
    }

    pub fn get_path_bits(self: *Population, id: u32, buffer: [*]u8, len: usize) void {
         for (self.paths) |p| {
             if (p.id == id) {
                 const copy_len = if (len < p.symbols.len) len else p.symbols.len;
                 @memcpy(buffer[0..copy_len], p.symbols[0..copy_len]);
                 return;
             }
         }
    }
};

// --- C Exports ---

var Gpa = std.heap.GeneralPurposeAllocator(.{}){};

// Returns an Opaque Pointer (Handle)
export fn ssn_create(cfg: *const Config) ?*Population {
    const allocator = Gpa.allocator();
    const ptr = Population.init(allocator, cfg.*) catch return null;
    return ptr;
}

export fn ssn_destroy(ptr: ?*Population) void {
    if (ptr) |p| {
        p.deinit();
    }
}

export fn ssn_select(ptr: ?*Population) u32 {
    if (ptr) |p| {
        return p.select();
    }
    return 0;
}

export fn ssn_update(ptr: ?*Population, selected_id: u32, reward_bits: u32) void {
    if (ptr) |p| {
        const reward: f32 = @bitCast(reward_bits);
        p.update(selected_id, reward);
    }
}

export fn ssn_get_path(ptr: ?*Population, id: u32, buffer: [*]u8, len: usize) void {
    if (ptr) |p| {
        p.get_path_bits(id, buffer, len);
    }
}
