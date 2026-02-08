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

<section class="flex flex-col gap-4">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl">{title}</h1>
        <a class="text-sm text-muted-foreground" href={url}>View all</a>
    </div>
    {@render children()}
</section>
