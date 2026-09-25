import { derived, writable } from "svelte/store";
import { ConfigService } from "$bindings/github.com/SladkyCitron/vlna/service";

export type Theme = "light" | "dark";

function isValidTheme(value: string): value is Theme {
  return value === "light" || value === "dark";
}

export const theme = writable<Theme>(undefined);

let configLoaded = false;
theme.subscribe((value) => {
  if (!configLoaded) {
    return;
  }

  ConfigService.SetTheme(value).catch((error) => {
    console.error("Failed to save theme:", error);
  });

  ConfigService.SaveConfig().catch((error) => {
    console.error("Failed to save config:", error);
  });
});

ConfigService.GetConfig()
  .then((config) => {
    if (config && isValidTheme(config.theme)) {
      theme.set(config.theme);
    }
    configLoaded = true;
  })
  .catch((error) => {
    configLoaded = true;
    console.error("Failed to load theme:", error);
  });

export const isDark = derived(theme, ($theme) => $theme === "dark");
