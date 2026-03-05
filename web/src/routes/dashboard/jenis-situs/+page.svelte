<script>
    // Mock data for site types
    const siteTypes = [
        { id: 1, name: 'Masjid', description: 'Tempat ibadah umat Islam', count: 24, color: 'primary' },
        { id: 2, name: 'Gereja', description: 'Tempat ibadah umat Kristen/Katolik', count: 8, color: 'secondary' },
        { id: 3, name: 'Pura', description: 'Tempat ibadah umat Hindu', count: 5, color: 'accent' },
        { id: 4, name: 'Vihara', description: 'Tempat ibadah umat Buddha', count: 3, color: 'info' },
        { id: 5, name: 'Klenteng', description: 'Tempat ibadah umat Khonghucu', count: 4, color: 'success' },
        { id: 6, name: 'Katedral', description: 'Gereja besar milik umat Katolik', count: 2, color: 'warning' },
        { id: 7, name: 'Surau', description: 'Tempat ibadah umat Islam (kecil)', count: 6, color: 'error' }
    ];
</script>

<div class="max-w-7xl mx-auto">
    <!-- Page Header -->
    <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">Jenis Situs Keagamaan</h1>
        <p class="text-gray-600">Kelola kategori jenis-jenis situs keagamaan yang terdaftar</p>
    </div>

    <!-- Add Type Button -->
    <div class="mb-6">
        <button class="btn btn-primary">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Tambah Jenis Situs
        </button>
    </div>

    <!-- Types Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each siteTypes as type (type.id)}
            <div class="card bg-base-100 shadow-xl border border-base-200">
                <div class="card-body">
                    <div class="flex items-start justify-between">
                        <div>
                            <h3 class="card-title text-lg font-bold">{type.name}</h3>
                            <p class="text-gray-600 text-sm">{type.description}</p>
                        </div>
                        <div class={`badge badge-${type.color}`}>
                            {type.count} situs
                        </div>
                    </div>

                    <div class="card-actions justify-end mt-4">
                        <button class="btn btn-outline btn-sm">Detail</button>
                        <button class="btn btn-primary btn-sm">Edit</button>
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div>
