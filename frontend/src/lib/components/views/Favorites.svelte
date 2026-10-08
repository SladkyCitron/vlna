<script lang="ts">
  import * as Item from "$lib/components/ui/item";
  import * as ButtonGroup from "$lib/components/ui/button-group";
  import { Button } from "$lib/components/ui/button";
  import { LoaderCircle, Radio, Heart, Play } from "@lucide/svelte";
  import { onMount } from "svelte";
  import type { Station } from "$bindings/github.com/SladkyCitron/vlna/service";
  import { openStationDetails } from "$lib/stores/stationDetails";
  import {
    favorites,
    favoritesLoading,
    loadFavorites,
    toggleFavorite,
  } from "$lib/stores/favorites";
  import { playStation } from "$lib/stores/player";
  import { getLocalizedCountryName } from "$lib/utils";
  import * as m from "$lib/paraglide/messages.js";

  function stopStationDetails(event: MouseEvent) {
    event.stopPropagation();
  }

  function removeFavorite(event: MouseEvent, station: Station) {
    stopStationDetails(event);
    toggleFavorite(station);
  }

  function play(event: MouseEvent, station: Station) {
    stopStationDetails(event);
    playStation(station);
  }

  onMount(() => {
    void loadFavorites();
  });
</script>

<div>
  <h1 class="pb-4 text-xl font-bold">{m.favorites()}</h1>
  {#if $favoritesLoading}
    <div class="flex items-center justify-center">
      <LoaderCircle class="h-8 w-8 animate-spin" />
      <p class="ml-2">{m.loading()}</p>
    </div>
  {:else}
    <div class="flex flex-col gap-4">
      {#each $favorites as station}
        <Item.Root variant="outline">
          {#snippet child({ props })}
            <a href="#/" onclick={() => openStationDetails(station)} {...props}>
              <Item.Media variant="image">
                {#if station.favicon}
                  <img
                    src={station.favicon}
                    alt={station.name}
                    class="size-16"
                  />
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
              <Item.Actions>
                <ButtonGroup.Root>
                  <Button
                    variant="outline"
                    onclick={(event) => removeFavorite(event, station)}
                  >
                    <Heart fill="currentColor" />
                  </Button>
                  <Button
                    variant="outline"
                    onclick={(event) => play(event, station)}
                  >
                    <Play />
                  </Button>
                </ButtonGroup.Root>
              </Item.Actions>
            </a>
          {/snippet}
        </Item.Root>
      {/each}
    </div>
  {/if}
</div>
