<script>
    import ContentCard from "$lib/components/top/ContentCard.svelte";
    import ContentGridWrapper from "$lib/components/top/ContentGridWrapper.svelte";

    let {
        albums
    } = $props()
</script>
<ContentGridWrapper title="albums" url="top/albums">
    <div class="grid grid-cols-3 gap-2">
        {#each albums as album, i (album.id)}
            <ContentCard title={album.title} pictureUrl={album.picture_url} contentType="albums"/>
        {/each}
    </div>
</ContentGridWrapper>
