<script>
	import { goto } from '$app/navigation';
	import apiService from '$lib/api';
	import { hasAllPermissions } from '$lib/permissions';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';
	import SitusTableActions from '$lib/components/ui/SitusTableActions.svelte';
	import { onMount } from 'svelte';

	let { data } = $props();
	let user = $derived(data.user);

	// Data states
	let sites = $state([]);
	let jenisSitusList = $state([]);
	let loading = $state(true);
	let error = $state(null);

	// Filter and search states
	let searchTerm = $state('');
	let statusFilter = $state('all');
	let typeFilter = $state('all');

	// Pagination states
	let currentPage = $state(1);
	let itemsPerPage = $state(10);

	// Modal and toast states
	let showDeleteModal = $state(false);
	let siteToDelete = $state(null);
	let isDeleting = $state(false);
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');

	onMount(async () => {
		await fetchData();
	});

	// Fetch sites and jenis situs on mount
	async function fetchData() {
		try {
			loading = true;
			error = null;

			// Fetch sites
			const sitesResponse = await apiService.get('/situs');
			const sitesData = await sitesResponse.json();
			sites = sitesData;

			// Fetch jenis situs for filter options
			const jenisResponse = await apiService.get('/jenis-situs');
			const jenisData = await jenisResponse.json();
			jenisSitusList = jenisData;
		} catch (err) {
			console.error('Failed to fetch data:', err);
			error = 'Gagal memuat data';
		} finally {
			loading = false;
		}
	}

	// Function to filter sites based on search and filters
	let filteredSites = $derived(
		sites.filter((site) => {
			const searchLower = searchTerm.toLowerCase();

			// Handle null/undefined values safely and add situs_id to search
			const siteIdKemenag = site.situs_id ? String(site.situs_id).toLowerCase() : '';
			const siteName = site.nama ? site.nama.toLowerCase() : '';
			const siteLocation = site.lokasi ? site.lokasi.toLowerCase() : '';
			const sitePendata = site.pendata ? site.pendata.toLowerCase() : '';

			const matchesSearch =
				siteName.includes(searchLower) ||
				siteLocation.includes(searchLower) ||
				sitePendata.includes(searchLower) ||
				siteIdKemenag.includes(searchLower);

			const matchesStatus = statusFilter === 'all' || site.status === statusFilter;
			const matchesType = typeFilter === 'all' || site.jenis === typeFilter;

			return matchesSearch && matchesStatus && matchesType;
		})
	);

	// Pagination Logic
	let totalPages = $derived(Math.ceil(filteredSites.length / itemsPerPage) || 1);

	let paginatedSites = $derived(
		filteredSites.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
	);

	// Reset ke halaman 1 setiap kali admin ngetik pencarian atau ganti filter
	$effect(() => {
		searchTerm;
		statusFilter;
		typeFilter;
		itemsPerPage;
		currentPage = 1;
	});

	function goToPage(page) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
		}
	}

	// Format date to Indonesian locale
	function formatDate(dateString) {
		if (!dateString) return '-';
		const date = new Date(dateString);
		return date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	// Get status badge class
	function getStatusBadgeClass(status) {
		switch (status) {
			case 'menunggu':
				return 'badge-warning';
			case 'terverifikasi':
				return 'badge-success';
			case 'ditolak':
				return 'badge-error';
			default:
				return 'badge-ghost';
		}
	}

	// Refresh data
	function handleRefresh() {
		fetchData();
	}

	// Confirm delete action
	function confirmDelete(site) {
		siteToDelete = site;
		showDeleteModal = true;
	}

	// Handle delete action
	async function handleDelete() {
		isDeleting = true;

		try {
			const res = await apiService.delete('/situs/' + siteToDelete.id);

			if (res.ok) {
				showDeleteModal = false;
				toastMessage = 'Data situs berhasil dihapus';
				toastType = 'success';
				showToast = true;
				fetchData();
				setTimeout(() => {
					showToast = false;
				}, 3000);
			} else {
				const errText = await res.text();
				toastMessage = 'Gagal menghapus data: ' + errText;
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			console.error('Delete error:', err);
			toastMessage = 'Terjadi kesalahan saat menghapus data';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			isDeleting = false;
		}
	}
</script>

<Toast
	show={showToast}
	message={toastMessage}
	type={toastType}
	onclose={() => (showToast = false)}
/>

<ConfirmModal
	show={showDeleteModal}
	title="Hapus Data Situs"
	message="Apakah Anda yakin ingin menghapus {siteToDelete?.nama}? Data yang dihapus tidak dapat dikembalikan."
	confirmText="Hapus"
	cancelText="Batal"
	isProcessing={isDeleting}
	onConfirm={handleDelete}
	onCancel={() => (showDeleteModal = false)}
/>

<div class="mx-auto max-w-7xl">
	<div class="mb-4 sm:mb-6">
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
			<div>
				<h1 class="text-2xl font-bold text-primary sm:text-3xl">Daftar Situs Keagamaan</h1>
				<p class="mt-1 text-sm text-base-content/70">
					Kelola dan pantau situs-situs keagamaan di wilayah Kecamatan Ciemas
				</p>
			</div>
			<button title="Refresh" onclick={handleRefresh} class="btn btn-square btn-ghost sm:btn-sm">
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
						d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
					/>
				</svg>
			</button>
		</div>
	</div>

	<div class="card mb-4 border border-base-200 bg-base-100 shadow-md sm:mb-6 sm:shadow-xl">
		<div class="card-body p-4 sm:p-5">
			<div class="grid grid-cols-1 gap-3 sm:gap-4 md:grid-cols-4">
				<div class="md:col-span-2">
					<div class="form-control">
						<input
							type="text"
							placeholder="Cari ID Kemenag, situs, lokasi, atau pendata..."
							class="input-bordered input min-h-11 w-full"
							bind:value={searchTerm}
						/>
					</div>
				</div>

				<div>
					<select class="select-bordered select min-h-11 w-full" bind:value={statusFilter}>
						<option value="all">Semua Status</option>
						<option value="menunggu">Menunggu</option>
						<option value="terverifikasi">Terverifikasi</option>
						<option value="ditolak">Ditolak</option>
					</select>
				</div>

				<div>
					<select class="select-bordered select min-h-11 w-full" bind:value={typeFilter}>
						<option value="all">Semua Jenis</option>
						{#each jenisSitusList as jenis (jenis.id)}
							<option value={jenis.nama}>{jenis.nama}</option>
						{/each}
					</select>
				</div>
			</div>

			<div class="mt-4 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
				<div>
					<span class="text-sm font-medium text-base-content/70">
						Menampilkan
						<span class="text-base-content">
							{filteredSites.length === 0 ? 0 : (currentPage - 1) * itemsPerPage + 1}
						</span>
						-
						<span class="text-base-content">
							{Math.min(currentPage * itemsPerPage, filteredSites.length)}
						</span>
						dari
						<span class="text-base-content">{filteredSites.length}</span> data
					</span>
				</div>
				<div class="flex w-full gap-2 sm:w-auto">
					{#if hasAllPermissions(user.permissions, ['situs:create'])}
						<button
							onclick={() => {
								goto('situs/tambah');
							}}
							class="btn min-h-11 flex-1 btn-primary sm:flex-none"
						>
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
							Tambah Situs
						</button>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<div class="card border border-base-200 bg-base-100 shadow-md sm:shadow-xl">
		<div class="card-body overflow-x-auto p-0 sm:p-5">
			{#if loading}
				<div class="flex flex-col items-center justify-center py-12">
					<span class="loading loading-lg loading-spinner text-primary"></span>
					<p class="mt-4 text-base-content/70">Memuat data...</p>
				</div>
			{:else if error}
				<div class="flex flex-col items-center justify-center py-12">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mb-2 h-12 w-12 text-error"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
					<p class="mb-4 text-error">{error}</p>
					<button onclick={fetchData} class="btn btn-sm btn-primary">Coba Lagi</button>
				</div>
			{:else if filteredSites.length === 0}
				<div class="flex flex-col items-center justify-center py-12">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mb-2 h-12 w-12 text-base-content/30"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
					<p class="text-base-content/70">Tidak ada situs yang ditemukan</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table w-full table-zebra">
						<thead>
							<tr>
								<th class="min-w-40">Nama Situs</th>
								<th class="min-w-30">ID (Kemenag)</th>
								<th class="min-w-30">Jenis</th>
								<th class="min-w-30">Lokasi</th>
								<th class="min-w-30">Pendata</th>
								<th class="min-w-25">Status Verifikasi</th>
								<th class="min-w-37.5">Terakhir Diubah</th>
								<th class="min-w-37.5">Aksi</th>
							</tr>
						</thead>
						<tbody>
							{#each paginatedSites as site (site.id)}
								<tr class="hover">
									<td>
										<div class="font-medium">{site.nama}</div>
									</td>
									<td>
										<span class="font-mono text-sm text-base-content/80">
											{site.situs_id || '-'}
										</span>
									</td>
									<td>
										<span class="badge badge-outline badge-sm">{site.jenis}</span>
									</td>
									<td>{site.lokasi}</td>
									<td>{site.pendata}</td>
									<td>
										<span class="badge badge-sm {getStatusBadgeClass(site.status)}">
											{site.status}
										</span>
									</td>
									<td class="text-sm text-base-content/70">{formatDate(site.updated_at)}</td>
									<td>
										<SitusTableActions
											id={site.id}
											status={site.status}
											permissions={data.user.permissions}
											onDelete={() => confirmDelete(site)}
										/>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				<div class="mt-6 mb-2 flex flex-col items-center justify-between gap-4 px-4 sm:flex-row">
					<div class="flex items-center gap-2" title="Pilih jumlah baris per halaman">
						<div
							class="flex h-8 w-8 items-center justify-center rounded-md bg-base-200 text-base-content/70"
						>
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
									d="M4 6h16M4 12h16M4 18h7"
								/>
							</svg>
						</div>
						<select class="select-bordered select select-sm font-medium" bind:value={itemsPerPage}>
							<option value={5}>5</option>
							<option value={10}>10</option>
							<option value={25}>25</option>
							<option value={50}>50</option>
						</select>
					</div>

					<div class="join">
						<button
							class="btn join-item px-2 btn-sm"
							disabled={currentPage === 1}
							onclick={() => goToPage(currentPage - 1)}
							title="Halaman Sebelumnya"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
							</svg>
						</button>

						<button
							class="no-animation btn pointer-events-none join-item flex items-center justify-center bg-base-200 px-4 font-medium text-base-content/80 btn-sm"
						>
							Halaman {currentPage} / {totalPages}
						</button>

						<button
							class="btn join-item px-2 btn-sm"
							disabled={currentPage === totalPages}
							onclick={() => goToPage(currentPage + 1)}
							title="Halaman Selanjutnya"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
							</svg>
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
