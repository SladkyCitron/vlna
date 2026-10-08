<script lang="ts">
  import * as m from "$lib/paraglide/messages.js";
  import { Separator } from "$lib/components/ui/separator";
  import * as Item from "$lib/components/ui/item";
  import * as Select from "$lib/components/ui/select";
  import { theme } from "$lib/stores/theme";
  import { locale } from "$lib/stores/locale";
  import { Palette, Languages } from "@lucide/svelte";
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
        <Select.Root type="single" bind:value={$theme}>
          <Select.Trigger class="w-45">
            <Select.Value placeholder={m["settings.theme"]()} />
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="light"
              >{m["settings.theme.light"]()}</Select.Item
            >
            <Select.Item value="dark">{m["settings.theme.dark"]()}</Select.Item>
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
        <Select.Root type="single" bind:value={$locale}>
          <Select.Trigger class="w-45">
            <Select.Value placeholder={m["settings.language"]()} />
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="en">English</Select.Item>
            <Select.Item value="sk">Slovenčina</Select.Item>
          </Select.Content>
        </Select.Root>
      </Item.Actions>
    </Item.Root>
  </Item.Group>
</div>
