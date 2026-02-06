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
        <h1 class="text-2xl">Top {title}</h1>
        <p class="text-sm text-muted-foreground">View all</p>
    </div>
    {@render children()}
</section>
