import { writable } from "svelte/store";
import { ConfigService } from "$bindings/github.com/SladkyCitron/vlna/service";
import { isLocale, setLocale } from "$lib/paraglide/runtime.js";

export const locale = writable<string>("en");

let configLoaded = false;
locale.subscribe((value) => {
  if (!configLoaded || !isLocale(value)) {
    return;
  }

  setLocale(value, { reload: false });
  ConfigService.SetLocale(value).catch((error) => {
    console.error("Failed to save locale:", error);
  });

  ConfigService.SaveConfig().catch((error) => {
    console.error("Failed to save config:", error);
  });
});

ConfigService.GetConfig()
  .then((config) => {
    const configuredLocale =
      config && "locale" in config && typeof config.locale === "string"
        ? config.locale.split("-")[0]
        : undefined;

    configLoaded = true;
    if (configuredLocale && isLocale(configuredLocale)) {
      locale.set(configuredLocale);
    }
  })
  .catch((error) => {
    configLoaded = true;
    console.error("Failed to load locale:", error);
  });
