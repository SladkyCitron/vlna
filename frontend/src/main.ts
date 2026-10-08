import "./app.css";
import { mount } from "svelte";
import App from "./App.svelte";
import { setLocale, isLocale } from "$lib/paraglide/runtime.js";
import { ConfigService } from "$bindings/github.com/SladkyCitron/vlna/service";

async function start() {
  try {
    const config = await ConfigService.GetConfig();
    const configuredLocale =
      config && "locale" in config && typeof config.locale === "string"
        ? config.locale.split("-")[0]
        : undefined;

    if (configuredLocale && isLocale(configuredLocale)) {
      console.log("Using locale:", configuredLocale);
      setLocale(configuredLocale, { reload: false });
    } else {
      console.log("Configured locale not found, using fallback English");
      setLocale("en", { reload: false });
    }
  } catch (error) {
    console.error(
      "Failed to load locale from config, using fallback English:",
      error,
    );
    setLocale("en", { reload: false });
  }

  mount(App, { target: document.getElementById("app")! });
}

start();
