<script>
	import apiService from '$lib/api';
	import { onMount } from 'svelte';
	import { hasAllPermissions } from '$lib/permissions';
	import Toast from '$lib/components/ui/Toast.svelte';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import RoleBadge from '$lib/components/ui/RoleBadge.svelte';
	import AddRoleForm from '$lib/components/ui/AddRoleForm.svelte';

	let { data } = $props();
	let user = $derived(data.user);
	let isAdding = $state(false);

	// Data statis
	const resources = [
		{ name: 'user', actions: ['create', 'read', 'update', 'delete'] },
		{
			name: 'situs',
			actions: ['create', 'read_all', 'read_own', 'update', 'delete', 'verify']
		},
		{
			name: 'jenis-situs',
			actions: ['create', 'read', 'update', 'delete']
		},
		{ name: 'role', actions: ['create', 'read', 'delete'] },
		{ name: 'policy', actions: ['create', 'read', 'update', 'delete'] }
	];

	// Helper untuk mendapatkan label yang lebih ramah
	function getActionLabel(action) {
		const labels = {
			create: 'Create',
			read: 'Read',
			read_all: 'Read All',
			read_own: 'Read Own',
			update: 'Update',
			update_all: 'Update All',
			update_own: 'Update Own',
			delete: 'Delete',
			verify: 'Verify'
		};
		return labels[action] || action;
	}

	// State reaktif Svelte 5 Runes
	let roles = $state([]);
	let policies = $state([]);
	let loading = $state(false);
	let updatingCells = $state(new Set());
	let addingRole = $state(false);

	// Modal and toast states
	let showDeleteModal = $state(false);
	let roleToDelete = $state(null);
	let isDeleting = $state(false);
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');

	// Derived state untuk kolom
	const columns = $derived(
		resources.flatMap((res) =>
			res.actions.map((act) => ({
				resource: res.name,
				action: act,
				label: `${capitalize(res.name)}: ${getActionLabel(act)}`
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

	// Wrapper untuk AddRoleForm
	function handleAddNewRoleWrapper(name) {
		addNewRole(name);
	}

	// Add new role
	async function addNewRole(name) {
		if (!name) return;

		addingRole = true;

		try {
			const res = await apiService.post('/roles', { name });
			if (res.ok) {
				toastMessage = 'Role berhasil ditambahkan';
				toastType = 'success';
				showToast = true;
				await loadRoles();
				setTimeout(() => {
					showToast = false;
				}, 3000);
			} else {
				const errorText = await res.text();
				toastMessage = errorText || 'Gagal menambahkan role';
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			toastMessage = 'Terjadi kesalahan saat menambahkan role';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			addingRole = false;
		}
	}

	// Confirm delete role
	function confirmDeleteRole(roleId, roleName) {
		roleToDelete = { id: roleId, name: roleName };
		showDeleteModal = true;
	}

	// Delete role
	async function deleteRole() {
		isDeleting = true;

		try {
			const res = await apiService.delete(`/roles/${roleToDelete.id}`);
			if (res.ok) {
				showDeleteModal = false;
				toastMessage = 'Role berhasil dihapus';
				toastType = 'success';
				showToast = true;
				await loadRoles();
				await loadPolicies();
				setTimeout(() => {
					showToast = false;
				}, 3000);
			} else {
				const errorText = await res.text();
				toastMessage = errorText || `Gagal menghapus role ${roleToDelete.name}`;
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			toastMessage = 'Terjadi kesalahan saat menghapus role';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			isDeleting = false;
			roleToDelete = null;
		}
	}

	// Load policies dari API
	async function loadPolicies() {
		loading = true;
		try {
			const res = await apiService.get('/policies');
			if (res.ok) {
				policies = await res.json();
			} else {
				const errorText = await res.text();
				toastMessage = errorText || 'Gagal memuat data policy';
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			toastMessage = 'Terjadi kesalahan saat memuat data';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			loading = false;
		}
	}

	// Add policy
	async function addPolicy(role, resource, action) {
		const key = getPolicyKey(role, resource, action);
		updatingCells.add(key);

		try {
			const res = await apiService.post('/policies', {
				subject: role,
				object: resource,
				action: action
			});

			if (res.ok) {
				policies.push({ subject: role, object: resource, action: action });
			} else {
				const errorText = await res.text();
				toastMessage = errorText || `Gagal menambahkan hak akses ${role} - ${resource}:${action}`;
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			toastMessage = 'Terjadi kesalahan saat menambahkan hak akses';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			updatingCells.delete(key);
		}
	}

	// Remove policy
	async function removePolicy(role, resource, action) {
		const key = getPolicyKey(role, resource, action);
		updatingCells.add(key);

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
				const errorText = await res.text();
				toastMessage = errorText || `Gagal mencabut hak akses ${role} - ${resource}:${action}`;
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			toastMessage = 'Terjadi kesalahan saat menghapus hak akses';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
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

<!-- Toast Notification -->
<Toast
	show={showToast}
	message={toastMessage}
	type={toastType}
	onclose={() => (showToast = false)}
/>

<!-- Confirmation Modal -->
<ConfirmModal
	show={showDeleteModal}
	title="Hapus Role"
	message="Apakah Anda yakin ingin menghapus role '{roleToDelete?.name}'? Semua hak akses terkait juga akan dihapus permanen."
	confirmText="Hapus"
	cancelText="Batal"
	isProcessing={isDeleting}
	onConfirm={deleteRole}
	onCancel={() => (showDeleteModal = false)}
/>

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
		<span> Centang untuk memberikan hak akses, hilangkan centang untuk mencabut hak akses. </span>
	</div>

	<!-- Add Role Form Component -->
	<AddRoleForm {isAdding} {user} onSubmit={handleAddNewRoleWrapper} />

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
					<table class="table table-zebra">
						<thead>
							<tr>
								<th class="bg-base-200 font-bold">Modul / Sumber Daya</th>
								{#each roles as role}
									<th class="bg-base-200 text-center">
										<div class="flex items-center justify-center gap-1">
											<RoleBadge name={role.name} showText={true} />
											{#if hasAllPermissions(user.permissions, 'role:delete')}
												<button
													class="btn text-error btn-ghost btn-xs hover:text-error"
													onclick={() => confirmDeleteRole(role.ID, role.name)}
													title="Hapus"
												>
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
												</button>
											{/if}
										</div>
									</th>
								{/each}
							</tr>
						</thead>
						<tbody>
							{#each resources as resource}
								<tr class="hover">
									<td class="font-bold capitalize">{resource.name}</td>
									{#each roles as role}
										<td>
											<div class="grid grid-cols-2 gap-2">
												{#each resource.actions as action}
													{@const isChecked = hasPolicy(role.name, resource.name, action)}
													{@const isUpdating = isCellUpdating(role.name, resource.name, action)}
													<label class="flex cursor-pointer items-center gap-2 text-xs">
														<input
															type="checkbox"
															class="checkbox checkbox-xs"
															checked={isChecked}
															disabled={isUpdating}
															onchange={(e) =>
																handleCheckboxChange(
																	role.name,
																	resource.name,
																	action,
																	e.target.checked
																)}
														/>
														<span class="text-base-content/80">
															{getActionLabel(action)}
														</span>
														{#if isUpdating}
															<span class="loading loading-xs loading-spinner text-primary"></span>
														{/if}
													</label>
												{/each}
											</div>
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
						<RoleBadge name={role.name} showText={true} />
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
