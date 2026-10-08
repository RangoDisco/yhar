<script lang="ts">
	import dayjs from 'dayjs';
	import * as Chart from '#lib/components/ui/chart/index.js';
	import { Chart as LayerChart, Calendar, Layer, Rect, Tooltip } from 'layerchart';
	import { scaleThreshold } from 'd3-scale';

	const lastDayOfYear = dayjs().toDate();
	const firstDayOfYear = dayjs(lastDayOfYear).subtract(1, 'year').toDate();
	const { data } = $props();

	const chartConfig = {
		background: {
			color: 'var(--card)'
		}
	} satisfies Chart.ChartConfig;
</script>

<Chart.Container config={chartConfig} class="aspect-6/1 min-h-38">
	<LayerChart
		{data}
		x="date"
		c="value"
		cScale={scaleThreshold()}
		cRange={['var(--chart-1-100)', 'var(--chart-1-300)', 'var(--chart-1-500)', 'var(--chart-1)']}
		cDomain={[25, 50, 75]}
		padding={{ top: 20, bottom: 0 }}
	>
		{#snippet children({ context })}
			<Layer type="html">
				<Calendar start={firstDayOfYear} end={lastDayOfYear}>
					{#snippet children({ cells, cellSize })}
						{#each cells as cell}
							{@const padding = 1}
							<Rect
								x={cell.x + padding}
								y={cell.y + padding}
								width={cellSize[0] - padding * 2}
								height={cellSize[1] - padding * 2}
								rx={4}
								fill={cell.color ?? 'var(--card)'}
								onpointermove={(e) => context.tooltip?.show(e, cell.data)}
								onpointerleave={(e) => context.tooltip?.hide()}
							/>
						{/each}
					{/snippet}
				</Calendar>
			</Layer>

			<Tooltip.Root>
				{#snippet children({ data })}
					<Tooltip.Header value={data.date} format={(d: Date) => dayjs(d).format('YYYY-MM-DD')} />

					{#if data.value != null}
						<Tooltip.List>
							<Tooltip.Item
								label="Scrobbles"
								value={data.value}
								format="integer"
								valueAlign="right"
							/>
						</Tooltip.List>
					{/if}
				{/snippet}
			</Tooltip.Root>
		{/snippet}
	</LayerChart>
</Chart.Container>
