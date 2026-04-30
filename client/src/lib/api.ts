import { browser } from "$app/environment";
import { goto } from "$app/navigation";

// ---------- Types ----------

export interface User {
  id: number;
  username: string;
  first_name: string;
  telegram_id?: string;
  fd_balance: number;
  created_at: string;
}

export interface FlowerTemplate {
  id: number;
  season: number;
  image_path: string;
  created_at: string;
}

export interface UserFlower {
  id: number;
  user_id: number;
  flower_id: number;
  day: number;
  last_watered_at?: string;
  is_dried: boolean;
  created_at: string;
}

export interface Seed {
  id: number;
  user_id: number;
  flower_id: number;
  quantity: number;
}

export interface Notification {
  id: number;
  user_id: number;
  type: string;
  payload: unknown;
  is_read: boolean;
  created_at: string;
}

export interface MeResponse {
  user: User;
  flowers: UserFlower[];
  seeds: Seed[];
}

export interface LoginResponse {
  token: string;
  user: User;
}

// ---------- Token ----------

export function getToken(): string | null {
  if (!browser) return null;
  return localStorage.getItem("token");
}

export function setToken(token: string): void {
  localStorage.setItem("token", token);
}

export function clearToken(): void {
  localStorage.removeItem("token");
}

// ---------- Core fetch ----------

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init.headers as Record<string, string>),
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(path, { ...init, headers });

  if (res.status === 401) {
    clearToken();
    goto("/auth");
    throw new ApiError(401, "Unauthorized");
  }

  if (!res.ok) {
    let msg = `HTTP ${res.status}`;
    try {
      const body = await res.json();
      msg = body.error || body.message || msg;
    } catch {
      // ignore
    }
    throw new ApiError(res.status, msg);
  }

  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

// ---------- Auth ----------

export const auth = {
  register(
    username: string,
    password: string,
    first_name: string,
  ): Promise<User> {
    return request("/api/auth/register", {
      method: "POST",
      body: JSON.stringify({ username, password, first_name }),
    });
  },
  login(username: string, password: string): Promise<LoginResponse> {
    return request("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    });
  },
};

// ---------- Me ----------

export const me = {
  get(): Promise<MeResponse> {
    return request("/api/me");
  },
  update(first_name: string): Promise<User> {
    return request("/api/me", {
      method: "PUT",
      body: JSON.stringify({ first_name }),
    });
  },
};

// ---------- Flowers ----------

export const flowers = {
  listTemplates(): Promise<FlowerTemplate[]> {
    return request("/api/flowers");
  },
  getTemplate(id: number): Promise<FlowerTemplate> {
    return request(`/api/flowers/${id}`);
  },
  getByUser(userId: number): Promise<UserFlower[]> {
    return request(`/api/flowers/user/${userId}`);
  },
  plant(flower_id: number): Promise<UserFlower> {
    return request("/api/flowers/plant", {
      method: "POST",
      body: JSON.stringify({ flower_id }),
    });
  },
  water(userFlowerId: number): Promise<{ status: string }> {
    return request(`/api/flowers/${userFlowerId}/water`, { method: "POST" });
  },
};

// ---------- Seeds ----------

export const seeds = {
  list(): Promise<Seed[]> {
    return request("/api/seeds");
  },
  share(
    to_user_id: number,
    flower_id: number,
    quantity: number,
  ): Promise<{ status: string }> {
    return request("/api/seeds/share", {
      method: "POST",
      body: JSON.stringify({ to_user_id, flower_id, quantity }),
    });
  },
};

// ---------- Leaderboard ----------

export const leaderboard = {
  get(limit = 100, offset = 0): Promise<User[]> {
    return request(`/api/leaderboard?limit=${limit}&offset=${offset}`);
  },
};

// ---------- Notifications ----------

export const notifications = {
  list(): Promise<Notification[]> {
    return request("/api/notifications");
  },
  markAllRead(): Promise<void> {
    return request("/api/notifications/read", { method: "PATCH" });
  },
};

// ---------- Users ----------

export const users = {
  getById(id: number): Promise<User> {
    return request(`/api/users/${id}`);
  },
  getByUsername(username: string): Promise<User> {
    return request(`/api/users/by-username/${encodeURIComponent(username)}`);
  },
};

// ---------- Dev ----------

export const dev = {
  tick(): Promise<{ status: string }> {
    return request("/api/dev/tick", { method: "POST" });
  },
  listUsers(): Promise<User[]> {
    return request("/api/dev/users");
  },
  giveSeeds(flower_id = 1, quantity = 5): Promise<{ status: string }> {
    return request("/api/dev/seeds", {
      method: "POST",
      body: JSON.stringify({ flower_id, quantity }),
    });
  },
  token(user_id: number): Promise<LoginResponse> {
    return request("/api/dev/token", {
      method: "POST",
      body: JSON.stringify({ user_id }),
    });
  },
  reset(): Promise<{ status: string }> {
    return request("/api/dev/reset", { method: "POST" });
  },
};

// ---------- Helpers ----------

export function totalFd(day: number): number {
  return (day * (day + 1)) / 2;
}

export function needsWatering(flower: UserFlower): boolean {
  if (flower.is_dried) return false;
  if (!flower.last_watered_at) return true;
  const today = new Date();
  const todayStr = `${today.getUTCFullYear()}-${String(today.getUTCMonth() + 1).padStart(2, "0")}-${String(today.getUTCDate()).padStart(2, "0")}`;
  return flower.last_watered_at.slice(0, 10) < todayStr;
}

export function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${String(d.getDate()).padStart(2, "0")}.${String(d.getMonth() + 1).padStart(2, "0")}.${d.getFullYear()}`;
}

// Flower emoji placeholders until real images are served
const FLOWER_EMOJIS: Record<number, string> = {
  1: "🌹",
  2: "🌷",
  3: "🌼",
};

export function flowerEmoji(flowerTemplateId: number): string {
  return FLOWER_EMOJIS[flowerTemplateId] ?? "🌸";
}
