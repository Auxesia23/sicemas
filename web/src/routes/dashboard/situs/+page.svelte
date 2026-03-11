<script>
	import { goto } from '$app/navigation';
	import apiService from '$lib/api';

	// Data states
	let sites = $state([]);
	let loading = $state(true);
	let error = $state(null);

	// Filter and search states
	let searchTerm = $state('');
	let statusFilter = $state('all');
	let typeFilter = $state('all');

	// Fetch sites on mount
	async function fetchSites() {
		try {
			loading = true;
			error = null;
			const response = await apiService.get('/situs');
			const data = await response.json();
			sites = data;
		} catch (err) {
			console.error('Failed to fetch sites:', err);
			error = 'Gagal memuat data situs';
		} finally {
			loading = false;
		}
	}

	// Initial fetch
	fetchSites();

	// Function to filter sites based on search and filters
	let filteredSites = $derived(
		sites.filter((site) => {
			const matchesSearch =
				site.nama.toLowerCase().includes(searchTerm.toLowerCase()) ||
				site.lokasi.toLowerCase().includes(searchTerm.toLowerCase()) ||
				site.pendata.toLowerCase().includes(searchTerm.toLowerCase());

			const matchesStatus = statusFilter === 'all' || site.status === statusFilter;
			const matchesType = typeFilter === 'all' || site.jenis === typeFilter;

			return matchesSearch && matchesStatus && matchesType;
		})
	);

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
			case 'disetujui':
				return 'badge-success';
			case 'ditolak':
				return 'badge-error';
			case 'revisi':
				return 'badge-info';
			default:
				return 'badge-ghost';
		}
	}

	// Refresh data
	function handleRefresh() {
		fetchSites();
	}
</script>

<div class="mx-auto max-w-7xl">
	<!-- Page Header -->
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

	<!-- Controls and Filters -->
	<div class="card mb-4 border border-base-200 bg-base-100 shadow-md sm:mb-6 sm:shadow-xl">
		<div class="card-body p-4 sm:p-5">
			<div class="grid grid-cols-1 gap-3 sm:gap-4 md:grid-cols-4">
				<div class="md:col-span-2">
					<div class="form-control">
						<input
							type="text"
							placeholder="Cari situs, lokasi, atau pendata..."
							class="input-bordered input min-h-11 w-full"
							bind:value={searchTerm}
						/>
					</div>
				</div>

				<div>
					<select class="select-bordered select min-h-11 w-full" bind:value={statusFilter}>
						<option value="all">Semua Status</option>
						<option value="menunggu">Menunggu</option>
						<option value="disetujui">Disetujui</option>
						<option value="ditolak">Ditolak</option>
						<option value="revisi">Revisi</option>
					</select>
				</div>

				<div>
					<select class="select-bordered select min-h-11 w-full" bind:value={typeFilter}>
						<option value="all">Semua Jenis</option>
						<option value="Masjid">Masjid</option>
						<option value="Gereja">Gereja</option>
						<option value="Pura">Pura</option>
						<option value="Vihara">Vihara</option>
						<option value="Klenteng">Klenteng</option>
					</select>
				</div>
			</div>

			<div class="mt-4 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
				<div>
					<span class="text-sm text-base-content/70">
						Menampilkan {filteredSites.length} dari {sites.length} situs
					</span>
				</div>
				<div class="flex w-full gap-2 sm:w-auto">
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
				</div>
			</div>
		</div>
	</div>

	<!-- Sites Table -->
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
					<button onclick={fetchSites} class="btn btn-sm btn-primary">Coba Lagi</button>
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
								<th class="min-w-50">Nama Situs</th>
								<th class="min-w-30">Jenis</th>
								<th class="min-w-30">Lokasi</th>
								<th class="min-w-30">Pendata</th>
								<th class="min-w-25">Status Verifikasi</th>
								<th class="min-w-37.5">Terakhir Diubah</th>
								<th class="min-w-37.5">Aksi</th>
							</tr>
						</thead>
						<tbody>
							{#each filteredSites as site (site.id)}
								<tr class="hover">
									<td>
										<div class="font-medium">{site.nama}</div>
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
									<td class="text-sm text-base-content/70">{formatDate(site.terakhir)}</td>
									<td>
										<div class="flex flex-wrap gap-1 sm:gap-2">
											<button class="btn btn-outline btn-xs sm:btn-sm">Lihat</button>
											<button class="btn btn-outline btn-xs btn-primary sm:btn-sm"> Edit </button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>
</div>
