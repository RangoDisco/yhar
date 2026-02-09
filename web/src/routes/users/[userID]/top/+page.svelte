<script>
    import * as Tabs from "$lib/components/ui/tabs/index";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
    import ContentCard from "$lib/components/top/ContentCard.svelte";

    const {data} = $props();
    const periods = ['week', 'month', 'year', 'overall'];
</script>
<div class="flex flex-col gap-8">
    <Tabs.Root value="week">
        <Tabs.List class="w-full">
            {#each periods as period}
                <Tabs.Trigger value={period}>{period.charAt(0).toUpperCase() + period.slice(1)}</Tabs.Trigger>
            {/each}
        </Tabs.List>
        {#each periods as period}
            <Tabs.Content value={period} class="flex flex-col gap-8">
                {#await data[period]}
                    Loading...
                {:then periodData}
                    <!--TOP ARTISTS-->
                    <ContentListWrapper title="Top artists" url="artists">
                        <div class="flex flex-wrap gap-4">
                            {#each periodData.artists.result as artist, i (artist.id)}
                                <ContentCard title={artist.name} pictureUrl={artist.picture_url} contentType="artists"
                                             contentID={artist.id}/>
                            {/each}
                        </div>
                    </ContentListWrapper>

                    <!--TOP ALBUMS-->
                    <ContentListWrapper title="Top albums" url="albums">
                        <div class="flex flex-wrap gap-4">
                            {#each periodData.albums.result as album, i (album.id)}
                                <ContentCard title={album.title} pictureUrl={album.picture_url}
                                             contentType="albums"/>
                            {/each}
                        </div>
                    </ContentListWrapper>

                    <!--TOP TRACKS-->
                    <ContentListWrapper title="Top tracks" url="tracks">
                        <div class="flex flex-col gap-4">
                            {#each periodData.tracks.result as track, i (track.id)}
                                <ContentListItem index={i} title={track.title}
                                                 pictureUrl={track.picture_url}
                                                 scrobbleCount={track.scrobble_count}
                                                 parentType="artists"
                                                 parents={track.artists}
                                                 contentType="tracks"/>
                            {/each}
                        </div>
                    </ContentListWrapper>
                {/await}
            </Tabs.Content>
        {/each}
    </Tabs.Root>

    <!--HISTORY-->
    <ContentListWrapper title="History" url="history">
        <div class="flex flex-col gap-2">
            {#each data.history.result as track}
                <HistoryListItem track={track} parentType="artists"/>
            {/each}
        </div>
    </ContentListWrapper>
</div>
