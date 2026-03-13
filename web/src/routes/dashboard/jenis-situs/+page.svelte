<script>
	import { onMount } from 'svelte';
	import apiService from '$lib/api';
	import { derived } from 'svelte/store';
	import { hasAllPermissions, hasAnyPermission } from '$lib/permissions';

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
	let formError = $state(null);
	let formSuccess = $state(null);

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
		formError = null;
		formSuccess = null;
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
		formError = null;
		formSuccess = null;

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

			formSuccess = `Jenis situs berhasil ${isEditMode ? 'diperbarui' : 'ditambahkan'}`;

			// Refresh data setelah berhasil
			await loadData();

			setTimeout(() => {
				closeModal();
			}, 1500);
		} catch (e) {
			formError = e.message || `Gagal ${isEditMode ? 'mengubah' : 'menambah'} jenis situs`;
		} finally {
			formLoading = false;
		}
	}

	async function deleteItem(id, nama) {
		if (!confirm(`Apakah Anda yakin ingin menghapus jenis situs "${nama}"?`)) {
			return;
		}

		try {
			const response = await apiService.delete(`/jenis-situs/${id}`);
			if (!response.ok) {
				const text = await response.text();
				throw new Error(text || 'Gagal menghapus jenis situs');
			}

			// Hapus item dari array secara reaktif
			siteTypes = siteTypes.filter((type) => type.id !== id);
		} catch (e) {
			alert(e.message || 'Terjadi kesalahan saat menghapus data');
		}
	}
</script>

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
									onclick={() => deleteItem(type.id, type.name)}
								>
									Hapus
								</button>
							{/if}
							{#if hasAllPermissions(user.permissions, ['jenis-situs:update'])}
								<button class="btn btn-sm btn-primary" onclick={() => openModal(type)}>
									Edit
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

		{#if formSuccess}
			<div class="mb-4 alert alert-success">
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
						d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
				<span>{formSuccess}</span>
			</div>
		{/if}

		{#if formError}
			<div class="mb-4 alert alert-error">
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
				<span>{formError}</span>
			</div>
		{/if}

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
