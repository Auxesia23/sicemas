<script>
	import { hasAllPermissions } from '$lib/permissions';
	import EditIcon from '$lib/components/icons/EditIcon.svelte';
	import ViewIcon from '$lib/components/icons/ViewIcon.svelte';
	import DeleteIcon from '$lib/components/icons/DeleteIcon.svelte';

	let { id, status, permissions, onDelete, basePath = '/dashboard/situs' } = $props();
</script>

<div class="flex flex-wrap gap-1 sm:gap-2">
	<!-- View Button -->
	<a
		href="{basePath}/{id}"
		class="btn btn-square btn-outline btn-xs sm:btn-sm"
		title="Lihat Detail"
	>
		<ViewIcon />
	</a>

	<!-- Edit Button (only if not verified) -->
	{#if hasAllPermissions(permissions, ['situs:update']) && status !== 'terverifikasi'}
		<a
			href="{basePath}/{id}/edit"
			class="btn btn-square btn-outline btn-xs btn-primary sm:btn-sm"
			title="Edit"
		>
			<EditIcon />
		</a>
	{/if}

	<!-- Delete Button -->
	{#if hasAllPermissions(permissions, ['situs:delete'])}
		<button
			class="btn btn-square btn-outline btn-xs btn-error sm:btn-sm"
			title="Hapus"
			onclick={onDelete}
		>
			<DeleteIcon />
		</button>
	{/if}
</div>
