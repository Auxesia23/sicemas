<script>
	import { hasAllPermissions } from '$lib/permissions';
	import EditIcon from '$lib/components/icons/EditIcon.svelte';
	import ViewIcon from '$lib/components/icons/ViewIcon.svelte';
	import DeleteIcon from '$lib/components/icons/DeleteIcon.svelte';

	let { id, permissions, onDelete, onUpdate } = $props();
</script>

<div class="flex flex-row flex-nowrap items-center justify-center gap-2">
	<!-- Edit Button -->
	{#if hasAllPermissions(permissions, ['user:update'])}
		<button
			class="btn btn-square btn-outline btn-xs btn-primary sm:btn-sm"
			title="Edit"
			onclick={onUpdate}
		>
			<EditIcon />
		</button>
	{/if}

	<!-- Delete Button -->
	{#if hasAllPermissions(permissions, ['user:delete'])}
		<button
			class="btn btn-square btn-outline btn-xs btn-error sm:btn-sm"
			title="Hapus"
			onclick={onDelete}
		>
			<DeleteIcon />
		</button>
	{/if}
</div>
