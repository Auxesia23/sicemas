<script>
	// Svelte 5 Runes: Receive detail prop from parent
	let { detail } = $props();
</script>

<!-- Detail Mushola View -->
<div class="space-y-6">
	<!-- Kalibrasi Arah Kiblat -->
	{#if detail.kalibrasi_arah_kiblat}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kalibrasi Arah Kiblat</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Azimut</p>
						<p class="font-mono font-medium">
							{detail.kalibrasi_arah_kiblat.azimut || '-'}
						</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Tanggal Kalibrasi</p>
						<p class="font-medium">
							{detail.kalibrasi_arah_kiblat.tanggal_kalibrasi || '-'}
						</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- SDM Mushola -->
	{#if detail.sdm_masjid}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">SDM Mushola</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Imam</div>
						<div class="stat-value text-primary">{detail.sdm_masjid.jumlah_imam || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Muadzin</div>
						<div class="stat-value text-secondary">{detail.sdm_masjid.jumlah_muadzin || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Jemaah</div>
						<div class="stat-value text-success">{detail.sdm_masjid.jumlah_jemaah || 0}</div>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Fasilitas Umum -->
	{#if detail.fasilitas_umum && detail.fasilitas_umum.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Fasilitas Umum</h2>
				<div class="flex flex-wrap gap-2">
					{#each detail.fasilitas_umum as fasilitas (fasilitas)}
						<span class="badge badge-secondary">{fasilitas}</span>
					{/each}
				</div>
			</div>
		</div>
	{/if}

	<!-- Nama Imam -->
	{#if detail.nama_imam && detail.nama_imam.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Nama Imam</h2>
				<div class="flex flex-wrap gap-2">
					{#each detail.nama_imam as imam (imam)}
						<span class="badge badge-outline">{imam}</span>
					{/each}
				</div>
			</div>
		</div>
	{/if}

	<!-- Nama Muadzin -->
	{#if detail.nama_muadzin && detail.nama_muadzin.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Nama Muadzin</h2>
				<div class="flex flex-wrap gap-2">
					{#each detail.nama_muadzin as muadzin (muadzin)}
						<span class="badge badge-outline">{muadzin}</span>
					{/each}
				</div>
			</div>
		</div>
	{/if}

	<!-- Kegiatan -->
	{#if detail.kegiatan && detail.kegiatan.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kegiatan</h2>
				<div class="flex flex-wrap gap-2">
					{#each detail.kegiatan as kegiatan (kegiatan)}
						<span class="badge badge-success">{kegiatan}</span>
					{/each}
				</div>
			</div>
		</div>
	{/if}
</div>
