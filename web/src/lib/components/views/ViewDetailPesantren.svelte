<script>
	// Svelte 5 Runes: Receive detail prop from parent
	let { detail } = $props();
</script>

<!-- Detail Pesantren View -->
<div class="space-y-6">
	<!-- Informasi Yayasan -->
	{#if detail.nama_yayasan || detail.pimpinan_pesantren}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Informasi Yayasan</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Nama Yayasan</p>
						<p class="font-medium">{detail.nama_yayasan || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Pimpinan Pesantren</p>
						<p class="font-medium">{detail.pimpinan_pesantren || '-'}</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Kepengurusan -->
	{#if detail.kepengurusan}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kepengurusan</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Ketua</p>
						<p class="font-medium">{detail.kepengurusan.ketua || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Sekretaris</p>
						<p class="font-medium">{detail.kepengurusan.sekretaris || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Bendahara</p>
						<p class="font-medium">{detail.kepengurusan.bendahara || '-'}</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Data Santri -->
	{#if detail.santri}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Data Santri</h2>
				<div class="overflow-x-auto">
					<table class="table table-zebra">
						<thead>
							<tr>
								<th>Kategori</th>
								<th class="text-center">Pria</th>
								<th class="text-center">Wanita</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td class="font-medium">Mondok</td>
								<td class="text-center">{detail.santri.mondok?.pria ?? 0}</td>
								<td class="text-center">{detail.santri.mondok?.wanita ?? 0}</td>
							</tr>
							<tr>
								<td class="font-medium">Tidak Mondok</td>
								<td class="text-center">{detail.santri.tidak_mondok?.pria ?? 0}</td>
								<td class="text-center">{detail.santri.tidak_mondok?.wanita ?? 0}</td>
							</tr>
							<tr class="font-bold">
								<td class="font-medium">Total</td>
								<td class="text-center">{detail.santri.total?.pria ?? 0}</td>
								<td class="text-center">{detail.santri.total?.wanita ?? 0}</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	{/if}

	<!-- Fasilitas Pondok -->
	{#if detail.fasilitas_pondok && detail.fasilitas_pondok.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Fasilitas Pondok</h2>
				<div class="flex flex-wrap gap-2">
					{#each detail.fasilitas_pondok as fasilitas (fasilitas)}
						<span class="badge badge-primary">{fasilitas}</span>
					{/each}
				</div>
			</div>
		</div>
	{/if}
</div>
