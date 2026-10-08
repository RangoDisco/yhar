<script lang="ts">
	import type { Snippet } from 'svelte';

	import { page } from '$app/state';
	import {cn} from "#lib/utils.ts";

	type Props = {
		title: string;
		url?: string | null;
		children: Snippet;
		class?: string | null;
	};

	let { title, url = $bindable(null), children, class: className }: Props = $props();
</script>

<section class={cn("flex w-full flex-col gap-3", className)}>
	<div class="flex items-center justify-between">
		<h1 class="text-sm md:text-base text-foreground/70">
			{#if url}
				<a
					class="font-medium text-foreground/70 hover:underline"
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
