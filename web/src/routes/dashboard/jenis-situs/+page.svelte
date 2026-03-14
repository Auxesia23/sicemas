<script>
	import { onMount } from 'svelte';
	import apiService from '$lib/api';
	import { hasAllPermissions } from '$lib/permissions';
	import DeleteIcon from '$lib/components/icons/DeleteIcon.svelte';
	import EditIcon from '$lib/components/icons/EditIcon.svelte';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';

	let { data } = $props();
	let user = $derived(data.user);

	// Menggunakan rune $state untuk Svelte 5
	let siteTypes = $state([]);
	let loading = $state(true);
	let error = $state(null);

	// Modal & Form state
	let modalOpen = $state(false);
	let isEditMode = $state(false);
	let editId = $state(null);
	let formNama = $state('');
	let formDeskripsi = $state('');
	let formLoading = $state(false);

	// Delete & Toast state
	let showDeleteModal = $state(false);
	let itemToDelete = $state(null);
	let isDeleting = $state(false);
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');

	const colorMap = {
		0: 'primary',
		1: 'secondary',
		2: 'accent',
		3: 'info',
		4: 'success',
		5: 'warning',
		6: 'error'
	};

	async function loadData() {
		loading = true;
		try {
			const response = await apiService.get('/jenis-situs');
			const data = await response.json();
			siteTypes = data.map((item, index) => ({
				id: item.id,
				name: item.nama,
				description: item.deskripsi,
				count: item.jumlah_situs,
				color: colorMap[index % 7]
			}));
			error = null;
		} catch (e) {
			error = 'Gagal memuat data jenis situs';
			console.error(e);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadData();
	});

	function openModal(item = null) {
		modalOpen = true;
		formLoading = false;

		if (item) {
			// Mode Edit
			isEditMode = true;
			editId = item.id;
			formNama = item.name;
			formDeskripsi = item.description;
		} else {
			// Mode Tambah
			isEditMode = false;
			editId = null;
			formNama = '';
			formDeskripsi = '';
		}
	}

	function closeModal() {
		modalOpen = false;
	}

	async function submitForm() {
		formLoading = true;

		try {
			const payload = {
				nama: formNama,
				deskripsi: formDeskripsi
			};

			let response;
			if (isEditMode) {
				// Proses Edit (PUT)
				response = await apiService.put(`/jenis-situs/${editId}`, payload);
			} else {
				// Proses Tambah (POST)
				response = await apiService.post('/jenis-situs', payload);
			}

			if (!response.ok) {
				const text = await response.text();
				throw new Error(text || `Gagal ${isEditMode ? 'mengubah' : 'menambah'} jenis situs`);
			}

			// Show success toast
			toastMessage = `Jenis situs berhasil ${isEditMode ? 'diperbarui' : 'ditambahkan'}`;
			toastType = 'success';
			showToast = true;

			// Refresh data setelah berhasil
			await loadData();

			closeModal();
			setTimeout(() => {
				showToast = false;
			}, 2000);
		} catch (e) {
			// Show error toast
			toastMessage = e.message || `Gagal ${isEditMode ? 'mengubah' : 'menambah'} jenis situs`;
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			formLoading = false;
		}
	}

	// Confirm delete action
	function confirmDeleteItem(id, name) {
		itemToDelete = { id, name };
		showDeleteModal = true;
	}

	// Handle delete action
	async function handleDelete() {
		isDeleting = true;

		try {
			const response = await apiService.delete(`/jenis-situs/${itemToDelete.id}`);
			if (!response.ok) {
				const text = await response.text();
				throw new Error(text || 'Gagal menghapus jenis situs');
			}

			// Success - show toast and update array
			toastMessage = 'Jenis situs berhasil dihapus';
			toastType = 'success';
			showToast = true;
			siteTypes = siteTypes.filter((type) => type.id !== itemToDelete.id);
			showDeleteModal = false;
			setTimeout(() => {
				showToast = false;
			}, 3000);
		} catch (e) {
			// Error - show error toast
			toastMessage = e.message || 'Terjadi kesalahan saat menghapus data';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			isDeleting = false;
			itemToDelete = null;
		}
	}
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
	title="Hapus Jenis Situs"
	message="Apakah Anda yakin ingin menghapus jenis situs '{itemToDelete?.name}'? Tindakan ini tidak dapat dibatalkan."
	confirmText="Hapus"
	cancelText="Batal"
	isProcessing={isDeleting}
	onConfirm={handleDelete}
	onCancel={() => (showDeleteModal = false)}
/>

<div class="mx-auto max-w-7xl">
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Jenis Situs Keagamaan</h1>
		<p class="text-gray-600">Kelola kategori jenis-jenis situs keagamaan yang terdaftar</p>
	</div>

	<div class="mb-6">
		{#if hasAllPermissions(user.permissions, ['jenis-situs:create'])}
			<button class="btn btn-primary" onclick={() => openModal()}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-1 h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 6v6m0 0v6m0-6h6m-6 0H6"
					/>
				</svg>
				Tambah Jenis Situs
			</button>
		{/if}
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<span class="loading loading-lg loading-spinner"></span>
		</div>
	{:else if error}
		<div class="alert alert-error">
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
		</div>
	{:else if siteTypes.length === 0}
		<div class="alert alert-info">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				class="h-6 w-6 shrink-0 stroke-current"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<span>Belum ada jenis situs keagamaan yang terdaftar</span>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each siteTypes as type (type.id)}
				<div class="card border border-base-200 bg-base-100 shadow-xl">
					<div class="card-body">
						<div class="flex items-start justify-between">
							<div>
								<h3 class="card-title text-lg font-bold">{type.name}</h3>
								<p class="text-sm text-gray-600">{type.description}</p>
							</div>
							<div class={`badge badge-${type.color}`}>
								{type.count} situs
							</div>
						</div>

						<div class="mt-4 card-actions justify-end">
							{#if hasAllPermissions(user.permissions, ['jenis-situs:delete'])}
								<button
									class="btn btn-outline btn-sm btn-error"
									onclick={() => confirmDeleteItem(type.id, type.name)}
									title="Hapus"
								>
									<DeleteIcon />
								</button>
							{/if}
							{#if hasAllPermissions(user.permissions, ['jenis-situs:update'])}
								<button
									class="btn btn-outline btn-sm btn-primary"
									onclick={() => openModal(type)}
									title="Edit"
								>
									<EditIcon />
								</button>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<dialog id="add_modal" class="modal" open={modalOpen}>
	<div class="modal-box">
		<form method="dialog">
			<button class="btn absolute top-2 right-2 btn-circle btn-ghost btn-sm" onclick={closeModal}>
				✕
			</button>
		</form>
		<h3 class="mb-4 text-lg font-bold">
			{isEditMode ? 'Edit Jenis Situs' : 'Tambah Jenis Situs'}
		</h3>

		<div class="form-control">
			<label class="label-text mb-2 block font-medium" for="formNama">Nama</label>
			<input
				id="formNama"
				type="text"
				placeholder="Masukkan nama jenis situs"
				class="input-bordered input w-full"
				bind:value={formNama}
				disabled={formLoading}
			/>
		</div>

		<div class="form-control mt-4">
			<label class="label-text mb-2 block font-medium" for="formDeskripsi">Deskripsi</label>
			<textarea
				id="formDeskripsi"
				class="textarea-bordered textarea h-24 w-full"
				placeholder="Masukkan deskripsi"
				bind:value={formDeskripsi}
				disabled={formLoading}
			></textarea>
		</div>

		<div class="modal-action">
			<button class="btn" onclick={closeModal} disabled={formLoading}>Batal</button>
			<button class="btn btn-primary" onclick={submitForm} disabled={formLoading}>
				{#if formLoading}
					<span class="loading loading-spinner"></span>
				{:else}
					Simpan
				{/if}
			</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button onclick={closeModal}>close</button>
	</form>
</dialog>
