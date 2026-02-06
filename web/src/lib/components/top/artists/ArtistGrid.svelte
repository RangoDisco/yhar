<script>
    import ContentCard from "$lib/components/top/ContentCard.svelte";
    import ContentGridWrapper from "$lib/components/top/ContentGridWrapper.svelte";

    let {
        artists
    } = $props()
</script>

<ContentGridWrapper title="artists" url="top/artists">

    <div class="grid grid-cols-3 gap-2">
        {#each artists as artist (artist.id)}
            <ContentCard title={artist.name} pictureUrl={artist.picture_url} contentType="artists"/>
        {/each}
    </div>
</ContentGridWrapper>
