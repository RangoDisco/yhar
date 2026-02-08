<script lang="ts">
    import ContentCard from "$lib/components/top/ContentCard.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";

    type Props = {
        items: [];
        title: string;
        url: string;
        contentType: "artists" | "albums";
    }

    let {
        items,
        title,
        url,
        contentType,
    }: Props = $props();
</script>

<ContentListWrapper title={title} url={url}>
    <div class="grid grid-cols-3 gap-2">
        {#each items as item, i (item.id)}
            <ContentCard title={item.title ?? item.name} pictureUrl={item.picture_url} contentType={contentType}/>
        {/each}
    </div>
</ContentListWrapper>
