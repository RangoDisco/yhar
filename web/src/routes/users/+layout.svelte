<script lang="ts">
	import { Button } from '#lib/components/ui/button/index.js';
	import ModeToggler from '#lib/components/ModeToggler.svelte';
	import Nav from '#lib/components/Nav.svelte';

	let { children, data } = $props();
</script>

<main class="flex flex-col">
	<div class="flex h-17 items-center gap-2 px-6">
		<Nav user={data.user} />
	</div>
	<hr />
	<div class="mx-auto w-full p-6 lg:max-w-450">
		{@render children()}
	</div>
</main>
