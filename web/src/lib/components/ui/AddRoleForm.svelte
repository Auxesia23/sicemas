<script>
	import { hasAllPermissions } from '$lib/permissions';

	let { isAdding = false, onSubmit, user } = $props();
	let roleName = $state('');

	function handleSubmit() {
		if (!roleName.trim()) return;
		onSubmit(roleName.trim());
		roleName = '';
	}
</script>

{#if hasAllPermissions(user.permissions, ['role:create'])}
	<div class="card mb-6 border border-base-200 bg-base-100 p-4 shadow">
		<h3 class="mb-3 text-lg font-bold">Tambah Role Baru</h3>
		<div class="flex gap-2">
			<input
				type="text"
				class="input-bordered input flex-1"
				placeholder="Nama role (contoh: developer)"
				bind:value={roleName}
				onkeydown={(e) => e.key === 'Enter' && handleSubmit()}
			/>
			<button
				class="btn btn-primary"
				onclick={handleSubmit}
				disabled={isAdding || !roleName.trim()}
			>
				{#if isAdding}
					<span class="loading loading-spinner"></span>
				{:else}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
				{/if}
				Tambah
			</button>
		</div>
	</div>
{/if}
