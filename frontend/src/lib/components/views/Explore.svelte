<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { MapPin, LoaderCircle } from "@lucide/svelte";
  import {
    IPInfoService,
    type IPInfo,
  } from "$bindings/github.com/SladkyCitron/vlna/service";
  import * as m from "$lib/paraglide/messages.js";
  import { onMount } from "svelte";

  let ipInfo: void | IPInfo | null = null;
  onMount(() => {
    IPInfoService.Fetch()
      .catch((err) => {
        console.error("Failed to fetch IP info:", err);
      })
      .then((info) => {
        ipInfo = info;
      });
  });

  function getLocalizedCountryName(countryCode: string): string {
    try {
      const displayNames = new Intl.DisplayNames([navigator.language], {
        type: "region",
      });
      return displayNames.of(countryCode) || countryCode;
    } catch (error) {
      console.error("Error getting localized country name:", error);
      return countryCode;
    }
  }
</script>

<div>
  <h1 class="pb-4 text-xl font-bold">{m.explore()}</h1>
  <Badge variant="secondary" class="flex items-center gap-2">
    {#if ipInfo}
      <MapPin class="h-4 w-4" />
      {ipInfo.cityName}, {getLocalizedCountryName(ipInfo.countryCode)}
    {:else}
      <LoaderCircle class="h-4 w-4 animate-spin" />
    {/if}
  </Badge>
</div>
