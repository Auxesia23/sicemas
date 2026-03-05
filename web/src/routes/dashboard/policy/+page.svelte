<script>
	import apiService from '$lib/api';
	import { onMount } from 'svelte';

	// Data statis
	const resources = ['user', 'situs', 'role', 'policy'];
	const actions = ['create', 'read', 'update', 'delete'];

	// State reaktif Svelte 5 Runes
	let roles = $state([]);
	let policies = $state([]);
	let loading = $state(false);
	let updatingCells = $state(new Set());
	let error = $state(null);
	let newRoleName = $state('');
	let addingRole = $state(false);
	let deletingRole = $state(null);

	// Derived state untuk kolom
	const columns = $derived(
		resources.flatMap((res) =>
			actions.map((act) => ({
				resource: res,
				action: act,
				label: `${capitalize(res)}: ${capitalize(act)}`
			}))
		)
	);

	// Helper function
	function capitalize(str) {
		return str.charAt(0).toUpperCase() + str.slice(1);
	}

	function getPolicyKey(role, resource, action) {
		return `${role}-${resource}-${action}`;
	}

	function hasPolicy(role, resource, action) {
		return policies.some((p) => p.subject === role && p.object === resource && p.action === action);
	}

	function isCellUpdating(role, resource, action) {
		return updatingCells.has(getPolicyKey(role, resource, action));
	}

	function getRoleBadgeClass(role) {
		if (role === 'admin') return 'badge-primary';
		if (role === 'operator') return 'badge-secondary';
		return 'badge-accent';
	}

	function getRoleBadgeIcon(role) {
		if (role === 'admin') {
			return 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z';
		}
		if (role === 'operator') {
			return 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z';
		}
		return 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z';
	}

	// Clear error after delay
	function clearError() {
		if (error) {
			setTimeout(() => {
				error = null;
			}, 5000);
		}
	}

	// Load roles dari API
	async function loadRoles() {
		try {
			const res = await apiService.get('/roles');
			if (res.ok) {
				const data = await res.json();
				roles = data;
			}
		} catch (err) {
			console.error('Gagal memuat roles:', err);
		}
	}

	// Add new role
	async function addNewRole() {
		if (!newRoleName.trim()) return;

		addingRole = true;
		error = null;

		try {
			const res = await apiService.post('/roles', { name: newRoleName.trim() });
			if (res.ok) {
				newRoleName = '';
				await loadRoles();
			} else {
				const errorText = await res.text();
				error = errorText || 'Gagal menambahkan role';
				clearError();
			}
		} catch (err) {
			error = 'Terjadi kesalahan saat menambahkan role';
			clearError();
		} finally {
			addingRole = false;
		}
	}

	// Delete role
	async function deleteRole(roleId, roleName) {
		if (!confirm(`Hapus role "${roleName}"? Semua hak akses terkait juga akan dihapus.`)) return;

		deletingRole = roleId;
		error = null;

		try {
			const res = await apiService.delete(`/roles/${roleId}`);
			if (res.ok) {
				await loadRoles();
				// Reload policies untuk menghapus policy yang terkait
				await loadPolicies();
			} else {
				const errorText = await res.text();
				error = errorText || `Gagal menghapus role ${roleName}`;
				clearError();
			}
		} catch (err) {
			error = 'Terjadi kesalahan saat menghapus role';
			clearError();
		} finally {
			deletingRole = null;
		}
	}

	// Load policies dari API
	async function loadPolicies() {
		loading = true;
		error = null;
		try {
			const res = await apiService.get('/policies');
			if (res.ok) {
				policies = await res.json();
			} else {
				const errorText = await res.text();
				error = errorText || 'Gagal memuat data policy';
				clearError();
			}
		} catch (err) {
			error = 'Terjadi kesalahan saat memuat data';
			clearError();
		} finally {
			loading = false;
		}
	}

	// Add policy
	async function addPolicy(role, resource, action) {
		const key = getPolicyKey(role, resource, action);
		updatingCells.add(key);
		error = null;

		try {
			const res = await apiService.post('/policies', {
				subject: role,
				object: resource,
				action: action
			});

			if (res.ok) {
				// API hanya return status code, buat policy lokal
				policies.push({ subject: role, object: resource, action: action });
			} else {
				// Handle error response (forbidden, etc.) - rollback checkbox
				const errorText = await res.text();
				error = errorText || `Gagal menambahkan hak akses ${role} - ${resource}:${action}`;
				clearError();
			}
		} catch (err) {
			// Error - rollback checkbox
			error = 'Terjadi kesalahan saat menambahkan hak akses';
			clearError();
		} finally {
			updatingCells.delete(key);
		}
	}

	// Remove policy
	async function removePolicy(role, resource, action) {
		const key = getPolicyKey(role, resource, action);
		updatingCells.add(key);
		error = null;

		try {
			const res = await apiService.delete('/policies', {
				subject: role,
				object: resource,
				action: action
			});

			if (res.ok) {
				policies = policies.filter(
					(p) => !(p.subject === role && p.object === resource && p.action === action)
				);
			} else {
				// Handle error response (forbidden, etc.) - rollback checkbox
				const errorText = await res.text();
				error = errorText || `Gagal mencabut hak akses ${role} - ${resource}:${action}`;
				clearError();
			}
		} catch (err) {
			// Error - rollback checkbox
			error = 'Terjadi kesalahan saat menghapus hak akses';
			clearError();
		} finally {
			updatingCells.delete(key);
		}
	}

	// Handle checkbox change
	function handleCheckboxChange(role, resource, action, checked) {
		if (checked) {
			addPolicy(role, resource, action);
		} else {
			removePolicy(role, resource, action);
		}
	}

	// Load data saat komponen di-mount
	onMount(async () => {
		await Promise.all([loadRoles(), loadPolicies()]);
	});
</script>

<div class="mx-auto max-w-7xl">
	<!-- Page Header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Manajemen Hak Akses</h1>
		<p class="text-gray-600">Kelola hak akses (policy matrix) untuk setiap peran dan sumber daya</p>
	</div>

	<!-- Info Card -->
	<div class="mb-6 alert alert-info shadow">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="h-6 w-6 shrink-0 stroke-current"
			fill="none"
			viewBox="0 0 24 24"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
			/>
		</svg>
		<span>
			Centang checkbox untuk memberikan hak akses, hilangkan centang untuk mencabut hak akses.
		</span>
	</div>

	<!-- Add Role Form -->
	<div class="card mb-6 border border-base-200 bg-base-100 p-4 shadow">
		<h3 class="mb-3 text-lg font-bold">Tambah Role Baru</h3>
		<div class="flex gap-2">
			<input
				type="text"
				class="input-bordered input flex-1"
				placeholder="Nama role (contoh: developer)"
				value={newRoleName}
				oninput={(e) => (newRoleName = e.target.value)}
				onkeydown={(e) => e.key === 'Enter' && addNewRole()}
			/>
			<button
				class="btn btn-primary"
				onclick={addNewRole}
				disabled={addingRole || !newRoleName.trim()}
			>
				{#if addingRole}
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

	<!-- Error Alert -->
	{#if error}
		<div class="mb-6 alert alert-error shadow">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-6 w-6 shrink-0 stroke-current"
				fill="none"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<span>{error}</span>
			<button class="btn btn-ghost btn-sm" onclick={() => (error = null)}>✕</button>
		</div>
	{/if}

	<!-- Loading State -->
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<span class="loading loading-lg loading-spinner text-primary"></span>
			<span class="ml-4 text-gray-600">Memuat data hak akses...</span>
		</div>
	{:else}
		<!-- Policy Matrix Table -->
		<div class="card border border-base-200 bg-base-100 shadow-xl">
			<div class="card-body">
				<div class="overflow-x-auto">
					<table class="table">
						<thead>
							<tr>
								<th class="min-w-[150px] bg-base-200">Peran / Fitur</th>
								{#each columns as col}
									<th class="min-w-[120px] bg-base-200 text-center">
										<div class="flex flex-col items-center">
											<span class="text-sm font-semibold">{col.label}</span>
										</div>
									</th>
								{/each}
							</tr>
						</thead>
						<tbody>
							{#each roles as role, index}
								<tr class={index % 2 === 1 ? 'bg-base-200/50' : ''}>
									<td class="font-semibold capitalize">
										<div class="flex items-center justify-between gap-2">
											<div class="flex items-center gap-2">
												<div
													class="badge badge-sm"
													class:badge-primary={role.name === 'admin'}
													class:badge-secondary={role.name === 'operator'}
													class:badge-accent={role.name !== 'admin' && role.name !== 'operator'}
												>
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="h-3 w-3"
														fill="none"
														viewBox="0 0 24 24"
														stroke="currentColor"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d={getRoleBadgeIcon(role.name)}
														/>
													</svg>
												</div>
												<span class="capitalize">{role.name}</span>
											</div>
											<button
												class="btn text-error btn-ghost btn-xs hover:text-error"
												onclick={() => deleteRole(role.ID, role.name)}
												disabled={deletingRole === role.ID}
												title="Hapus role"
											>
												{#if deletingRole === role.ID}
													<span class="loading loading-xs loading-spinner"></span>
												{:else}
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="h-4 w-4"
														fill="none"
														viewBox="0 0 24 24"
														stroke="currentColor"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
														/>
													</svg>
												{/if}
											</button>
										</div>
									</td>
									{#each columns as col}
										{@const isChecked = hasPolicy(role.name, col.resource, col.action)}
										{@const isUpdating = isCellUpdating(role.name, col.resource, col.action)}
										<td class="text-center">
											<label class="flex cursor-pointer justify-center gap-2">
												<input
													type="checkbox"
													class="checkbox checkbox-sm"
													checked={isChecked}
													disabled={isUpdating}
													onchange={(e) =>
														handleCheckboxChange(
															role.name,
															col.resource,
															col.action,
															e.target.checked
														)}
												/>
												{#if isUpdating}
													<span class="loading loading-xs loading-spinner text-primary"></span>
												{/if}
											</label>
										</td>
									{/each}
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				<!-- Legend -->
				<div class="mt-6 flex flex-wrap gap-4">
					{#each roles as role}
						<div class="flex items-center gap-2">
							<div
								class="badge badge-sm"
								class:badge-primary={role.name === 'admin'}
								class:badge-secondary={role.name === 'operator'}
								class:badge-accent={role.name !== 'admin' && role.name !== 'operator'}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-3 w-3"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d={getRoleBadgeIcon(role.name)}
									/>
								</svg>
							</div>
							<span class="text-sm text-gray-600 capitalize">{role.name}</span>
						</div>
					{/each}
				</div>
			</div>
		</div>
	{/if}

	<!-- Summary Stats -->
	{#if !loading}
		<div class="mt-6 grid grid-cols-1 gap-4 md:grid-cols-3">
			<div class="stat rounded-box border border-base-200 bg-base-100 shadow">
				<div class="stat-title">Total Kebijakan</div>
				<div class="stat-value text-primary">{policies.length}</div>
			</div>
			<div class="stat rounded-box border border-base-200 bg-base-100 shadow">
				<div class="stat-title">Peran Aktif</div>
				<div class="stat-value text-secondary">{roles.length}</div>
			</div>
			<div class="stat rounded-box border border-base-200 bg-base-100 shadow">
				<div class="stat-title">Sumber Daya</div>
				<div class="stat-value text-accent">{resources.length}</div>
			</div>
		</div>
	{/if}
</div>
