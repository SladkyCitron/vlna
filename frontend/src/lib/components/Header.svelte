<script lang="ts">
  import { onMount } from "svelte";
  import { Window, Events } from "@wailsio/runtime";
  import { Button } from "$lib/components/ui/button";
  import { Minus, Square, Copy, X } from "@lucide/svelte";

  let maximized = $state(false);

  onMount(() => {
    Window.IsMaximised().then((isMaximized) => (maximized = isMaximized));

    const unoff = Events.On("wails:window:maximise", () => (maximized = true));
    const unoffRestore = Events.On(
      "wails:window:unmaximise",
      () => (maximized = false),
    );

    return () => {
      unoff?.();
      unoffRestore?.();
    };
  });

  async function toggleMaximise() {
    const isMaximized = await Window.IsMaximised();
    if (isMaximized) {
      await Window.UnMaximise();
    } else {
      await Window.Maximise();
    }
    maximized = !isMaximized;
  }
</script>

<header
  class="col-span-2 h-8 wails-drag border-b border-border flex items-center justify-between pl-4 bg-popover select-none"
>
  <div class="flex items-center gap-2 text-xs font-bold text-muted-foreground">
    Vlna
  </div>

  <div class="wails-no-drag flex h-full items-center">
    <Button
      variant="ghost"
      class="h-8 w-11 rounded-none hover:bg-accent focus-visible:ring-0"
      onclick={Window.Minimise}
    >
      <Minus class="h-3.5 w-3.5" strokeWidth={1.5} />
    </Button>

    <Button
      variant="ghost"
      class="h-8 w-11 rounded-none hover:bg-accent focus-visible:ring-0"
      onclick={toggleMaximise}
    >
      {#if maximized}
        <Copy class="h-3.5 w-3.5 rotate-90" strokeWidth={1.5} />
      {:else}
        <Square class="h-3.5 w-3.5" strokeWidth={1.5} />
      {/if}
    </Button>

    <Button
      variant="ghost"
      class="h-8 w-11 rounded-none hover:bg-destructive hover:text-destructive-foreground focus-visible:ring-0"
      onclick={Window.Close}
    >
      <X class="h-4 w-4" strokeWidth={1.5} />
    </Button>
  </div>
</header>
