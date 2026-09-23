import "./app.css";
import { mount } from "svelte";
import App from "./App.svelte";
import { setLocale, isLocale } from "$lib/paraglide/runtime.js";

// set locale
const systemLocale = navigator.language.split("-")[0];
if (isLocale(systemLocale)) {
  console.log("Using locale:", systemLocale);
  setLocale(systemLocale, { reload: false });
} else {
  console.log(`Locale ${systemLocale} not found, using fallback English`);
  setLocale("en", { reload: false });
}

mount(App, { target: document.getElementById("app")! });
