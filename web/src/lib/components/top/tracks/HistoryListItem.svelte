<script lang="ts">
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Trash, ListMusic } from '@lucide/svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import * as Tooltip from '#lib/components/ui/tooltip/index.js';

	import { page } from '$app/state';
	import type { Scrobble } from '#lib/types/content.js';

	dayjs.extend(relativeTime);

	type Props = {
		scrobble: Scrobble;
		parentType: 'artists' | 'albums';
		handleDelete: (id: string) => void;
	};

	let { scrobble, parentType, handleDelete }: Props = $props();
</script>

<article class="flex items-center justify-between">
	<div class="flex w-full min-w-0 items-center gap-4">
		<Avatar.Root class="h-8 w-8 rounded-md lg:h-10 lg:w-10">
			<Avatar.Image src={scrobble.track.picture_url} alt={`${scrobble.track.title}'s picture`} />
			<Avatar.Fallback class="h-8 w-8 rounded-md">
				<ListMusic size={18} class="text-muted-foreground" />
			</Avatar.Fallback>
		</Avatar.Root>
		<div class="flex w-full min-w-0 flex-col">
			<div class="flex min-w-0 items-center gap-2">
				<span class="line-clamp-1 w-full min-w-0 text-lg">{scrobble.track.title}</span>
				<Tooltip.Root>
					<Tooltip.Trigger>
						<p class="w-40 text-sm whitespace-nowrap text-muted-foreground">
							{dayjs(scrobble.scrobbled_at).fromNow()}
						</p>
					</Tooltip.Trigger>
					<Tooltip.Content>
						<p class="text-sm">{dayjs(scrobble.scrobbled_at).format('YYYY-MM-DD HH:mm ')}</p>
					</Tooltip.Content>
				</Tooltip.Root>
			</div>
			<div class="flex gap-1">
				{#if parentType === 'artists'}
					{#each scrobble.track.artists as artist, i}
						{#if i !== 0}
							·
						{/if}
						<a
							class="h-6 text-muted-foreground hover:underline"
							href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a
						>
					{/each}
				{:else}
					<a
						class="text-muted-foreground hover:underline"
						href="/users/{page.params.userID}/top/albums/{scrobble.track.album.id}"
						>{scrobble.track.album.title}</a
					>
				{/if}
			</div>
		</div>
	</div>
	<Button
		variant="outline"
		size="icon"
		onclick={() => handleDelete(scrobble.id)}
		aria-label="Delete scrobble"
	>
		<Trash />
	</Button>
</article>
