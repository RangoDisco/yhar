<script lang="ts">
	import ContentListWrapper from '#lib/components/top/ContentListWrapper.svelte';
	import ContentListItem from '#lib/components/top/ContentListItem.svelte';
	import { page } from '$app/state';
	import UploadableImage from '#lib/components/top/UploadableImage.svelte';
	import SvelteHead from '#lib/components/SvelteHead.svelte';
	import MusicBrainzLogo from '#lib/components/logo/MusicBrainzLogo.svelte';
	import { Button } from '#lib/components/ui/button/index.js';
	import HistoryList from "#lib/components/top/tracks/HistoryList.svelte";

	let { data } = $props();
</script>

<SvelteHead title={data.track.title} description="Top tracks for {data.track.title}" />
<main class="flex flex-col gap-8">
	<section class="flex items-center gap-4">
		<UploadableImage
			pictureUrl={`${data.track.pictureDomain}/images/${data.track.picturePath}`}
			alt="{data.track.title}'s picture"
			contentType="albums"
			contentID={parseInt(data.track.id)}
			uploadEnabled={data.user?.role === 'ADMIN'}
		/>
		<div class="flex flex-col">
			<p class="text-3xl lg:text-4xl">{data.track.title}</p>
			<div class="flex items-center gap-1">
				{#each data.track.artists as artist, i}
					{#if i < 3}
						{#if i !== 0}
							·
						{/if}
						<a
							class="text-sm text-muted-foreground hover:underline"
							href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a
						>
					{/if}
				{/each}
			</div>
			<p class="text-muted-foreground">{data.track.scrobble_count} scrobbles</p>
			{#if data.track.music_brainz_id != null}
				<Button
					variant="ghost"
					size="icon-lg"
					href="https://musicbrainz.org/release/{data.track.music_brainz_id}"
					target="_blank"
				>
					<MusicBrainzLogo />
				</Button>
			{/if}
		</div>
	</section>
	<HistoryList
			scrobbles={data.history.results}
			url="history?tracks={data.track.id}"
			parentType="albums"
	/>
<!--	<div class="flex flex-col gap-8">-->
<!--		&lt;!&ndash;TOP TRACKS&ndash;&gt;-->
<!--		<ContentListWrapper title="Top tracks">-->
<!--			<div class="flex flex-col gap-4">-->
<!--				{#each data.tracks.results as track, i (track.id)}-->
<!--					<ContentListItem-->
<!--						index={i + 1}-->
<!--						title={track.title}-->
<!--						pictureUrl={track.picture_url}-->
<!--						scrobbleCount={track.scrobble_count}-->
<!--						parentType="artists"-->
<!--						parents={track.artists}-->
<!--						contentType="tracks"-->
<!--					/>-->
<!--				{/each}-->
<!--			</div>-->
<!--		</ContentListWrapper>-->
<!--	</div>-->
</main>
