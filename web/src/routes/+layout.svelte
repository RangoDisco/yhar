<script lang="ts">
	import './layout.css';
	import { Toaster } from '#lib/components/ui/sonner/index.js';
	import { ModeWatcher } from 'mode-watcher';
	import { TooltipProvider } from '#lib/components/ui/tooltip/index.js';

	let { children } = $props();
</script>

<main class="min-h-screen flex flex-col md:flex-row">
	<ModeWatcher />
	<Toaster />
	<TooltipProvider>
		{@render children()}
	</TooltipProvider>
</main>
