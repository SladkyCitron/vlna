<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import * as Item from "$lib/components/ui/item";
  import { MapPin, LoaderCircle, Radio } from "@lucide/svelte";
  import {
    IPInfoService,
    type IPInfo,
    StationService,
    type Stations,
  } from "$bindings/github.com/SladkyCitron/vlna/service";
  import * as m from "$lib/paraglide/messages.js";
  import { onMount } from "svelte";

  // fetch IP info for location
  let ipInfo: IPInfo | null = null;

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

  // fetch stations
  let stations: Stations = [];
  onMount(() => {
    const fetchExploreData = async () => {
      try {
        ipInfo = await IPInfoService.Fetch();
        if (!ipInfo) {
          return;
        }

        stations = await StationService.GetStationsByCountryCode(
          ipInfo.countryCode
        );
      } catch (error) {
        console.error("Failed to fetch Explore data:", error);
      }
    };

    void fetchExploreData();
  });
</script>

<div>
  <h1 class="pb-4 text-xl font-bold">{m.explore()}</h1>
  <Badge variant="secondary" class="mb-4 flex items-center gap-2">
    {#if ipInfo}
      <MapPin class="h-4 w-4" />
      {ipInfo.cityName}, {getLocalizedCountryName(ipInfo.countryCode)}
    {:else}
      <LoaderCircle class="h-4 w-4 animate-spin" />
    {/if}
  </Badge>
  {#if stations}
    <div class="flex flex-col gap-4">
      {#each stations as station}
        <Item.Root variant="outline">
          <Item.Media variant="image">
            {#if station.favicon}
              <img src={station.favicon} alt={station.name} class="size-16" />
            {:else}
              <Radio class="size-8" aria-label={station.name} />
            {/if}
          </Item.Media>
          <Item.Content>
            <Item.Title>{station.name}</Item.Title>
            <Item.Description>
              {getLocalizedCountryName(station.countrycode)}
              {#if station.tags}
                | {station.tags
                  .split(",")
                  .slice(0, 4) // limit to first 5 tags only
                  .join(", ")}
              {/if}
            </Item.Description>
          </Item.Content>
        </Item.Root>
      {/each}
    </div>
  {:else}
    <div class="flex items-center justify-center">
      <LoaderCircle class="h-8 w-8 animate-spin" />
      <p class="ml-2">{m.loading()}</p>
    </div>
  {/if}
</div>
