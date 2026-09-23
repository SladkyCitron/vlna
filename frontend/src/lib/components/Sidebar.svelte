<script lang="ts">
  import { activeView, type View } from "$lib/stores/nav";
  import { Compass, Heart, Settings } from "@lucide/svelte";
  import * as m from "$lib/paraglide/messages.js";

  const navItems: { id: View; label: string; icon: any }[] = [
    { id: "explore", label: m.explore(), icon: Compass },
    { id: "favorites", label: m.favorites(), icon: Heart },
  ];

  function setActiveView(view: View) {
    activeView.set(view);
  }
</script>

<aside
  class="border-border flex w-60 shrink-0 flex-col justify-between border-r p-3"
>
  <nav class="space-y-1">
    <div
      class="text-muted-foreground px-4 py-2 text-sm font-semibold tracking-wider"
    >
      {m.library()}
    </div>
    {#each navItems as item}
      <button
        type="button"
        class="hover:bg-accent hover:text-accent-foreground text-foreground flex w-full items-center gap-2 rounded-md px-4 py-2 text-sm font-medium focus:outline-none"
        class:bg-accent={item.id === $activeView}
        class:text-accent-foreground={item.id === $activeView}
        onclick={() => setActiveView(item.id)}
      >
        <svelte:component this={item.icon} class="h-4 w-4" />
        {item.label}
      </button>
    {/each}
    <div class="border-border mt-2 border-t py-2">
      <button
        type="button"
        class="hover:bg-accent hover:text-accent-foreground text-foreground flex w-full items-center gap-2 rounded-md px-4 py-2 text-sm font-medium focus:outline-none"
        class:bg-accent={"settings" === $activeView}
        class:text-accent-foreground={"settings" === $activeView}
        onclick={() => setActiveView("settings")}
      >
        <Settings class="h-4 w-4" />
        {m.settings()}
      </button>
    </div>
  </nav>
</aside>
