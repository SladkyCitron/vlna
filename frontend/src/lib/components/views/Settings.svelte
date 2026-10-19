<script lang="ts">
  import * as m from "$lib/paraglide/messages.js";
  import { Separator } from "$lib/components/ui/separator";
  import * as Item from "$lib/components/ui/item";
  import * as Select from "$lib/components/ui/select";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import { Switch } from "$lib/components/ui/switch";
  import { theme } from "$lib/stores/theme";
  import { locale } from "$lib/stores/locale";
  import { Palette, Languages, Radio, Info } from "@lucide/svelte";
  import {
    ConfigService,
    VersionService,
    PlayerService,
  } from "$bindings/github.com/SladkyCitron/vlna/service";
  import { onMount } from "svelte";

  let version = $state("");
  let discordRPC = $state(true);

  onMount(async () => {
    const [loadedVersion, config] = await Promise.all([
      VersionService.GetVersion(),
      ConfigService.GetConfig(),
    ]);
    version = loadedVersion;

    if (
      config &&
      "discordrpc" in config &&
      typeof config.discordrpc === "boolean"
    ) {
      discordRPC = config.discordrpc;
    }
  });

  async function saveDiscordRPC(enabled: boolean) {
    try {
      await ConfigService.SetDiscordRPC(enabled);
      await ConfigService.SaveConfig();
      await PlayerService.UpdateDiscordRPC();
    } catch (error) {
      console.error("Failed to save Discord Rich Presence setting:", error);
    }
  }

  async function openGitHub() {
    try {
      await VersionService.OpenGitHub();
    } catch (error) {
      console.error("Failed to open GitHub:", error);
    }
  }

  const themes = [
    { value: "light", label: m["settings.theme.light"]() },
    { value: "dark", label: m["settings.theme.dark"]() },
  ];

  const locales = [
    { value: "en", label: "English" },
    { value: "en-pt", label: "English (Pirate Speak)" },
    { value: "sk", label: "Slovenčina" },
  ];
</script>

<div class="flex flex-col gap-4">
  <h1 class="pb-4 text-xl font-bold">{m.settings()}</h1>
  <Separator />
  <Item.Group>
    <Item.Root>
      <Item.Media variant="icon"><Palette /></Item.Media>
      <Item.Content>
        <Item.Title>{m["settings.theme"]()}</Item.Title>
        <Item.Description>{m["settings.theme.description"]()}</Item.Description>
      </Item.Content>
      <Item.Actions>
        <Select.Root type="single" bind:value={$theme} items={themes}>
          <Select.Trigger class="w-45">
            <Select.Value placeholder={m["settings.theme"]()} />
          </Select.Trigger>
          <Select.Content>
            {#each themes as themeOption}
              <Select.Item value={themeOption.value}>
                {themeOption.label}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </Item.Actions>
    </Item.Root>
    <Item.Root>
      <Item.Media variant="icon"><Languages /></Item.Media>
      <Item.Content>
        <Item.Title>{m["settings.language"]()}</Item.Title>
        <Item.Description
          >{m["settings.language.description"]()}</Item.Description
        >
      </Item.Content>
      <Item.Actions>
        <Select.Root type="single" bind:value={$locale} items={locales}>
          <Select.Trigger class="w-45">
            <Select.Value placeholder={m["settings.language"]()} />
          </Select.Trigger>
          <Select.Content>
            {#each locales as localeOption}
              <Select.Item value={localeOption.value}>
                {localeOption.label}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </Item.Actions>
    </Item.Root>
    <Item.Root>
      <Item.Media variant="icon"><Radio /></Item.Media>
      <Item.Content>
        <Item.Title>{m["settings.discordrpc"]()}</Item.Title>
        <Item.Description>
          {m["settings.discordrpc.description"]()}
        </Item.Description>
      </Item.Content>
      <Item.Actions>
        <Switch bind:checked={discordRPC} onCheckedChange={saveDiscordRPC} />
      </Item.Actions>
    </Item.Root>
    <Item.Root>
      <Item.Media variant="icon"><Info /></Item.Media>
      <Item.Content>
        <Item.Title>{m["settings.about"]()}</Item.Title>
      </Item.Content>
      <Item.Actions>
        <Dialog.Root>
          <Dialog.Trigger
            type="button"
            class={buttonVariants({ variant: "outline" })}
          >
            {m["settings.about"]()}
          </Dialog.Trigger>
          <Dialog.Content>
            <Dialog.Header>
              <Dialog.Title>{m["settings.about"]()}</Dialog.Title>
              <Dialog.Description>
                Vlna v{version}
                <br />
                (c) 2026 SladkyCitron; MIT License
              </Dialog.Description>
            </Dialog.Header>
            <Button onclick={openGitHub}>GitHub</Button>
          </Dialog.Content>
        </Dialog.Root>
      </Item.Actions>
    </Item.Root>
  </Item.Group>
</div>
