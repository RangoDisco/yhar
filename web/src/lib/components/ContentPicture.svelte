<script lang="ts">
    import type { ContentType } from "#lib/types/content.ts";
    import {ListMusic} from "@lucide/svelte";
    import Disc3 from "@lucide/svelte/icons/disc-3";
    import User from "@lucide/svelte/icons/user";
    import type { ClassValue } from "clsx";
    import { cn } from "#lib/utils.ts";

    type Props = {
        img?: HTMLImageElement | null;
        class?: ClassValue;
        url: string | null;
        alt: string;
        contentType: ContentType;
        restProps?: Record<string, unknown>;
    };

    let {
        img = $bindable(null),
        class: className,
        url,
        alt,
        contentType,
        ...restProps
    }: Props = $props();

    let imageError: boolean = $state(false);

    function getFallbackIcon() {
        switch (contentType) {
            case "artists":
                return User;
            case "albums":
                return Disc3;
            case "tracks":
                return ListMusic;
        }
    }

    const FallbackIcon = $derived(getFallbackIcon());
</script>

<div
        class={cn(
		"flex justify-center items-center aspect-square bg-muted overflow-hidden",
		contentType === "artists" ? "rounded-full" : "rounded-md",
		className
	)}
        {...restProps}
>
    {#if url && !imageError}
        <img
                bind:this={img}
                class="aspect-square w-full"
                src={url}
                {alt}
                onerror={() => (imageError = true)}
        />
    {:else}
        <div class="flex justify-center items-center p-3 w-full">
        <FallbackIcon class="text-muted-foreground w-full h-full" />
        </div>
    {/if}
</div>
