<script>
	// Mock data for sites
	const sites = [
		{
			id: 1,
			name: 'Masjid Al-Muttaqin',
			type: 'Masjid',
			location: 'Desa Ciemas',
			status: 'Aktif',
			lastUpdated: '2023-05-15',
			contact: 'Ahmad Fauzi'
		},
		{
			id: 2,
			name: 'Gereja Katolik Santo Yosef',
			type: 'Gereja',
			location: 'Desa Citorek',
			status: 'Aktif',
			lastUpdated: '2023-05-14',
			contact: 'Siti Nurhaliza'
		},
		{
			id: 3,
			name: 'Pura Tirta Suci',
			type: 'Pura',
			location: 'Desa Cibodas',
			status: 'Non-Aktif',
			lastUpdated: '2023-05-10',
			contact: 'Budi Santoso'
		},
		{
			id: 4,
			name: 'Vihara Dharma Bhakti',
			type: 'Vihara',
			location: 'Desa Cikahuripan',
			status: 'Aktif',
			lastUpdated: '2023-05-12',
			contact: 'Lina Marlina'
		},
		{
			id: 5,
			name: 'Klenteng Hok Tek Bio',
			type: 'Klenteng',
			location: 'Desa Surade',
			status: 'Verifikasi',
			lastUpdated: '2023-05-13',
			contact: 'Rizki Pratama'
		},
		{
			id: 6,
			name: 'Masjid Al-Huda',
			type: 'Masjid',
			location: 'Desa Cikadu',
			status: 'Aktif',
			lastUpdated: '2023-05-11',
			contact: 'Ahmad Fauzi'
		},
		{
			id: 7,
			name: 'Gereja Injili Indonesia',
			type: 'Gereja',
			location: 'Desa Cibolang',
			status: 'Aktif',
			lastUpdated: '2023-05-09',
			contact: 'Siti Nurhaliza'
		},
		{
			id: 8,
			name: 'Pura Lempuyang',
			type: 'Pura',
			location: 'Desa Cikiray',
			status: 'Non-Aktif',
			lastUpdated: '2023-05-08',
			contact: 'Budi Santoso'
		}
	];

	// Filter and search states
	let searchTerm = $state('');
	let statusFilter = $state('all');
	let typeFilter = $state('all');

	// Function to filter sites based on search and filters
	let filteredSites = $derived(
		sites.filter((site) => {
			const matchesSearch =
				site.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
				site.location.toLowerCase().includes(searchTerm.toLowerCase()) ||
				site.contact.toLowerCase().includes(searchTerm.toLowerCase());

			const matchesStatus = statusFilter === 'all' || site.status === statusFilter;
			const matchesType = typeFilter === 'all' || site.type === typeFilter;

			return matchesSearch && matchesStatus && matchesType;
		})
	);
</script>

<div class="mx-auto max-w-7xl">
	<!-- Page Header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Daftar Situs Keagamaan</h1>
		<p class="text-gray-600">Kelola dan pantau situs-situs keagamaan di wilayah Kecamatan Ciemas</p>
	</div>

	<!-- Controls and Filters -->
	<div class="card mb-6 border border-base-200 bg-base-100 shadow-xl">
		<div class="card-body">
			<div class="grid grid-cols-1 gap-4 md:grid-cols-4">
				<div class="md:col-span-2">
					<div class="form-control">
						<input
							type="text"
							placeholder="Cari situs, lokasi, atau kontak..."
							class="input-bordered input w-full"
							bind:value={searchTerm}
						/>
					</div>
				</div>

				<div>
					<select class="select-bordered select w-full" bind:value={statusFilter}>
						<option value="all">Semua Status</option>
						<option value="Aktif">Aktif</option>
						<option value="Non-Aktif">Non-Aktif</option>
						<option value="Verifikasi">Verifikasi</option>
					</select>
				</div>

				<div>
					<select class="select-bordered select w-full" bind:value={typeFilter}>
						<option value="all">Semua Jenis</option>
						<option value="Masjid">Masjid</option>
						<option value="Gereja">Gereja</option>
						<option value="Pura">Pura</option>
						<option value="Vihara">Vihara</option>
						<option value="Klenteng">Klenteng</option>
					</select>
				</div>
			</div>

			<div class="mt-4 flex items-center justify-between">
				<div>
					<span class="text-sm text-gray-500">
						Menampilkan {filteredSites.length} dari {sites.length} situs
					</span>
				</div>
				<div>
					<button class="btn btn-primary">
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
	<div class="card border border-base-200 bg-base-100 shadow-xl">
		<div class="card-body overflow-x-auto">
			<table class="table table-zebra">
				<thead>
					<tr>
						<th>Nama Situs</th>
						<th>Jenis</th>
						<th>Lokasi</th>
						<th>Kontak</th>
						<th>Status</th>
						<th>Terakhir Diubah</th>
						<th>Aksi</th>
					</tr>
				</thead>
				<tbody>
					{#each filteredSites as site (site.id)}
						<tr>
							<td>
								<div class="font-medium">{site.name}</div>
							</td>
							<td>
								<span class="badge badge-outline">{site.type}</span>
							</td>
							<td>{site.location}</td>
							<td>{site.contact}</td>
							<td>
								{#if site.status === 'Aktif'}
									<span class="badge badge-success">{site.status}</span>
								{:else if site.status === 'Non-Aktif'}
									<span class="badge badge-error">{site.status}</span>
								{:else}
									<span class="badge badge-warning">{site.status}</span>
								{/if}
							</td>
							<td>{site.lastUpdated}</td>
							<td>
								<div class="flex space-x-2">
									<button class="btn btn-outline btn-xs">Lihat</button>
									<button class="btn btn-outline btn-xs btn-primary">Edit</button>
								</div>
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="7" class="text-center py-8">
								<div class="flex flex-col items-center justify-center">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-12 w-12 text-gray-400 mb-2"
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
									<p class="text-gray-500">Tidak ada situs yang ditemukan</p>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
