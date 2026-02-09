<script lang="ts">
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import ContentCard from "$lib/components/top/ContentCard.svelte";

    let {data} = $props();
</script>

<main class="flex flex-col gap-8">
    <section class="flex items-center gap-4">
        <img class="rounded-full aspect-square h-24" src={data.artist.picture_url} alt="{data.artist.name}'s picture"/>
        <div class="flex flex-col">
            <p>{data.artist.name}</p>
            <p class="text-muted-foreground">{data.artist.scrobble_count} scrobbles</p>
        </div>
    </section>

    <div class="flex flex-col gap-8">
        <!--TOP ALBUMS-->
        <ContentListWrapper title="Top albums" url="albums?artist={data.artist.id}">
            <div class="flex flex-wrap gap-4">
                {#each data.albums.result as album, i (album.id)}
                    <ContentCard title={album.title ?? album.name} pictureUrl={album.picture_url} contentType="albums"
                                 contentID={album.id}/>
                {/each}
            </div>
        </ContentListWrapper>

        <!--TOP TRACKS-->
        <ContentListWrapper title="Top tracks" url="tracks?artist={data.artist.id}">
            <div class="flex flex-col gap-4">
                {#each data.tracks.result as track, i (track.id)}
                    <ContentListItem index={i} title={track.title}
                                     pictureUrl={track.picture_url}
                                     scrobbleCount={track.scrobble_count}
                                     parentType="albums"
                                     parents={[track.album]}
                                     contentType="tracks"/>
                {/each}
            </div>
        </ContentListWrapper>

        <!--HISTORY-->
        <ContentListWrapper title="History" url="history?artist={data.artist.id}">
            <div class="flex flex-col gap-2">
                {#each data.history.result as track}
                    <HistoryListItem track={track} parentType="albums"/>
                {/each}
            </div>
        </ContentListWrapper>
    </div>
</main>
