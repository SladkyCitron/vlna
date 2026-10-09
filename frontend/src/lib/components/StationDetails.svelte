<script lang="ts">
  import { Check, Radio, Copy } from "@lucide/svelte";
  import * as Sheet from "$lib/components/ui/sheet";
  import * as Item from "$lib/components/ui/item";
  import { Separator } from "$lib/components/ui/separator";
  import { Button } from "$lib/components/ui/button";
  import {
    selectedStation,
    isStationDetailsOpen,
    closeStationDetails,
  } from "$lib/stores/stationDetails";
  import { getLocalizedCountryName } from "$lib/utils";
  import * as m from "$lib/paraglide/messages.js";

  function handleClose(open: boolean) {
    if (!open) {
      open = false;
      setTimeout(() => {
        closeStationDetails();
      }, 200);
    }
  }

  let copiedField = $state<string | null>(null);
  let copiedResetTimeout: ReturnType<typeof setTimeout> | undefined;

  async function copy(text: string, field: string) {
    try {
      await navigator.clipboard.writeText(text);
      copiedField = field;
      if (copiedResetTimeout) {
        clearTimeout(copiedResetTimeout);
      }
      copiedResetTimeout = setTimeout(() => {
        if (copiedField === field) {
          copiedField = null;
        }
      }, 1000);
    } catch (err) {
      console.error("Failed to copy text: ", err);
    }
  }
</script>

<Sheet.Root bind:open={$isStationDetailsOpen} onOpenChange={handleClose}>
  <Sheet.Content side="right">
    <Sheet.Header>
      <Sheet.Title>{m.stationDetails()}</Sheet.Title>
    </Sheet.Header>
    {#if $selectedStation}
      <div class="flex flex-col overflow-x-hidden overflow-y-auto px-4">
        {#if $selectedStation.favicon}
          <img
            src={$selectedStation.favicon}
            alt={$selectedStation.name}
            class="border-border flex size-32 items-center justify-center self-center rounded-md border"
          />
        {:else}
          <div
            class="border-border flex items-center justify-center self-center rounded-md border p-4"
          >
            <Radio class="size-32" />
          </div>
        {/if}
        <p class="my-4 text-center text-2xl font-bold">
          {$selectedStation.name}
        </p>
        <Separator />
        <Item.Group>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.url()}</Item.Title>
              <Item.Description>{$selectedStation.url}</Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() => copy($selectedStation.url, "url")}
              >
                {#if copiedField === "url"}<Check class="size-4" />{:else}<Copy
                    class="size-4"
                  />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.homepage()}</Item.Title>
              <Item.Description>{$selectedStation.homepage}</Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() => copy($selectedStation.homepage, "homepage")}
              >
                {#if copiedField === "homepage"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.tags()}</Item.Title>
              <Item.Description>{$selectedStation.tags}</Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() => copy($selectedStation.tags, "tags")}
              >
                {#if copiedField === "tags"}<Check class="size-4" />{:else}<Copy
                    class="size-4"
                  />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.country()}</Item.Title>
              <Item.Description>
                {getLocalizedCountryName($selectedStation.countrycode)}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() =>
                  copy(
                    getLocalizedCountryName($selectedStation.countrycode),
                    "country"
                  )}
              >
                {#if copiedField === "country"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.countrycode()}</Item.Title>
              <Item.Description>
                {$selectedStation.countrycode}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() =>
                  copy($selectedStation.countrycode, "countrycode")}
              >
                {#if copiedField === "countrycode"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.state()}</Item.Title>
              <Item.Description>
                {$selectedStation.state}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() => copy($selectedStation.state, "state")}
              >
                {#if copiedField === "state"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.languages()}</Item.Title>
              <Item.Description>
                {$selectedStation.languagecodes
                  .split(",")
                  .map((code) => getLocalizedCountryName(code))
                  .join(", ")}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() =>
                  copy(
                    $selectedStation.languagecodes
                      .split(",")
                      .map((code) => getLocalizedCountryName(code))
                      .join(", "),
                    "languages"
                  )}
              >
                {#if copiedField === "languages"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.languagecodes()}</Item.Title>
              <Item.Description>
                {$selectedStation.languagecodes}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() =>
                  copy($selectedStation.languagecodes, "languagecodes")}
              >
                {#if copiedField === "languagecodes"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.votes()}</Item.Title>
              <Item.Description>
                {$selectedStation.votes}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() => copy($selectedStation.votes.toString(), "votes")}
              >
                {#if copiedField === "votes"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.codec()}</Item.Title>
              <Item.Description>
                {$selectedStation.codec}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() => copy($selectedStation.codec, "codec")}
              >
                {#if copiedField === "codec"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.bitrate()}</Item.Title>
              <Item.Description>
                {$selectedStation.bitrate} kbps
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() =>
                  copy(
                    $selectedStation.bitrate.toString() + " kbps",
                    "bitrate"
                  )}
              >
                {#if copiedField === "bitrate"}<Check
                    class="size-4"
                  />{:else}<Copy class="size-4" />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
          <Item.Root>
            <Item.Content>
              <Item.Title>{m.hls()}</Item.Title>
              <Item.Description>
                {#if $selectedStation.hls === 1}
                  {m.yes()}
                {:else}
                  {m.no()}
                {/if}
              </Item.Description>
            </Item.Content>
            <Item.Actions>
              <Button
                variant="secondary"
                onclick={() =>
                  copy($selectedStation.hls === 1 ? m.yes() : m.no(), "hls")}
              >
                {#if copiedField === "hls"}<Check class="size-4" />{:else}<Copy
                    class="size-4"
                  />{/if}
              </Button>
            </Item.Actions>
          </Item.Root>
        </Item.Group>
      </div>
    {/if}
  </Sheet.Content>
</Sheet.Root>

<style>
  Check,
  Copy {
    transition:
      background-color 0.2s ease-in-out,
      color 0.2s ease-in-out;
  }
</style>
