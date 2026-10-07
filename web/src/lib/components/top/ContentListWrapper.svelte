<script lang="ts">
	import type { Snippet } from 'svelte';

	import { page } from '$app/state';

	type Props = {
		title: string;
		url?: string | null;
		children: Snippet;
	};

	let { title, url = $bindable(null), children }: Props = $props();
</script>

<section class="flex w-full flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-lg text-foreground/90 md:text-xl">
			{#if url}
				<a
					class="text-lg text-foreground/90 hover:underline md:text-xl"
					href="/users/{page.params.userID}/{url}"
				>
					{title}
				</a>
			{:else}
				{title}
			{/if}
		</h1>
	</div>
	{@render children()}
</section>
