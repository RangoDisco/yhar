<script lang="ts">
    import type {Snippet} from "svelte";

    type Props = {
        title: string;
        url: string;
        children: Snippet
    }

    let {
        title,
        url,
        children
    }: Props = $props();
</script>

<section class="flex flex-col gap-6">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl md:text-3xl">{title}</h1>
        <a class="text-sm text-muted-foreground md:text-base hover:underline" href={url}>View all</a>
    </div>
    {@render children()}
</section>
