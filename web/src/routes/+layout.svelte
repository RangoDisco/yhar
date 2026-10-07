<script lang="ts">
	import './layout.css';
	import { Toaster } from '#lib/components/ui/sonner/index.js';
	import { ModeWatcher } from 'mode-watcher';
	import { TooltipProvider } from '#lib/components/ui/tooltip/index.js';

	let { children } = $props();
</script>

<section class="min-h-screen">
	<ModeWatcher />
	<Toaster />
	<TooltipProvider>
		{@render children()}
	</TooltipProvider>
</section>
