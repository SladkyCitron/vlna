import { writable } from "svelte/store";

export type View = "explore" | "favorites" | "settings";

export const activeView = writable<View>("explore");
