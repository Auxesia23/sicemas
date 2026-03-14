<script>
	let { show = false, message = '', type = 'success', onclose } = $props();
</script>

{#if show}
	<div class="toast toast-end toast-top z-9999">
		<div
			class="alert shadow-lg {type === 'success'
				? 'alert-success text-white'
				: 'alert-error text-white'}"
		>
			<div>
				<h3 class="font-bold">{type === 'success' ? 'Berhasil!' : 'Gagal!'}</h3>
				<div class="text-sm">{message}</div>
			</div>
			{#if onclose}
				<button class="btn btn-circle btn-ghost btn-sm" onclick={onclose}>✕</button>
			{/if}
		</div>
	</div>
{/if}
