<script>
    // Mock data for statistics
    const siteDistribution = [
        { name: 'Masjid', value: 24, percentage: 46, color: 'bg-primary' },
        { name: 'Gereja', value: 8, percentage: 15, color: 'bg-secondary' },
        { name: 'Pura', value: 5, percentage: 10, color: 'bg-accent' },
        { name: 'Vihara', value: 3, percentage: 6, color: 'bg-info' },
        { name: 'Klenteng', value: 4, percentage: 8, color: 'bg-success' },
        { name: 'Lainnya', value: 8, percentage: 15, color: 'bg-warning' }
    ];

    const monthlyData = [
        { month: 'Jan', sites: 48, visitors: 1200 },
        { month: 'Feb', sites: 49, visitors: 1350 },
        { month: 'Mar', sites: 50, visitors: 1420 },
        { month: 'Apr', sites: 51, visitors: 1580 },
        { month: 'May', sites: 52, visitors: 1650 },
        { month: 'Jun', sites: 52, visitors: 1720 }
    ];

    const statusData = [
        { status: 'Aktif', count: 42, percentage: 80, color: 'bg-success' },
        { status: 'Verifikasi', count: 5, percentage: 10, color: 'bg-warning' },
        { status: 'Non-Aktif', count: 5, percentage: 10, color: 'bg-error' }
    ];
</script>

<div class="max-w-7xl mx-auto">
    <!-- Page Header -->
    <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">Statistik dan Analisis</h1>
        <p class="text-gray-600">Analisis data dan statistik terkait situs keagamaan di wilayah Kecamatan Ciemas</p>
    </div>

    <!-- Stats Summary -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body">
                <h3 class="text-3xl font-bold text-primary">52</h3>
                <p class="text-gray-500">Total Situs</p>
                <div class="mt-2">
                    <div class="flex items-center text-sm text-success">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                        </svg>
                        <span>+2.5% dari bulan lalu</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body">
                <h3 class="text-3xl font-bold text-secondary">12</h3>
                <p class="text-gray-500">Jenis Situs</p>
                <div class="mt-2">
                    <div class="flex items-center text-sm text-success">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                        </svg>
                        <span>+1 jenis baru</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body">
                <h3 class="text-3xl font-bold text-accent">42</h3>
                <p class="text-gray-500">Situs Aktif</p>
                <div class="mt-2">
                    <div class="flex items-center text-sm text-success">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                        </svg>
                        <span>+5 dari bulan lalu</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body">
                <h3 class="text-3xl font-bold text-info">98%</h3>
                <p class="text-gray-500">Akurasi Data</p>
                <div class="mt-2">
                    <div class="flex items-center text-sm text-success">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                        </svg>
                        <span>+1.2% dari bulan lalu</span>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Charts Section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <!-- Site Distribution Chart -->
        <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body">
                <h2 class="card-title text-lg font-bold mb-4">Distribusi Situs Berdasarkan Jenis</h2>
                <div class="space-y-4">
                    {#each siteDistribution as item (item.name)}
                        <div>
                            <div class="flex justify-between mb-1">
                                <span class="text-sm font-medium">{item.name}</span>
                                <span class="text-sm font-medium">{item.value} situs ({item.percentage}%)</span>
                            </div>
                            <div class="w-full bg-base-200 rounded-full h-2.5">
                                <div class={`${item.color} h-2.5 rounded-full`} style="width: {item.percentage}%"></div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        </div>

        <!-- Monthly Trend Chart -->
        <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body">
                <h2 class="card-title text-lg font-bold mb-4">Tren Perkembangan Situs (6 Bulan Terakhir)</h2>
                <div class="h-64 flex items-end justify-between pt-8 pb-4">
                    {#each monthlyData as data (data.month)}
                        <div class="flex flex-col items-center flex-1">
                            <div class="text-xs mb-2">{data.sites}</div>
                            <div
                                class="w-8 bg-primary rounded-t-lg"
                                style="height: {data.sites * 3}px; min-height: 10px;"
                            ></div>
                            <div class="text-xs mt-2">{data.month}</div>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>

    <!-- Status Distribution -->
    <div class="card bg-base-100 shadow-xl border border-base-200 mb-8">
        <div class="card-body">
            <h2 class="card-title text-lg font-bold mb-4">Status Situs Saat Ini</h2>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                {#each statusData as item (item.status)}
                    <div class="border border-base-200 rounded-lg p-4">
                        <div class="flex items-center justify-between mb-2">
                            <h3 class="font-bold">{item.status}</h3>
                            <span class="text-lg font-bold">{item.count}</span>
                        </div>
                        <div class="w-full bg-base-200 rounded-full h-2.5">
                            <div class={`${item.color} h-2.5 rounded-full`} style="width: {item.percentage}%"></div>
                        </div>
                        <div class="text-right text-sm mt-1">{item.percentage}%</div>
                    </div>
                {/each}
            </div>
        </div>
    </div>

    <!-- Export Options -->
    <div class="card bg-base-100 shadow-xl border border-base-200">
        <div class="card-body">
            <h2 class="card-title text-lg font-bold mb-4">Ekspor Data Statistik</h2>
            <div class="flex flex-wrap gap-4">
                <button class="btn btn-outline">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                    Ekspor PDF
                </button>
                <button class="btn btn-outline">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                    Ekspor Excel
                </button>
                <button class="btn btn-outline">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                    Ekspor CSV
                </button>
            </div>
        </div>
    </div>
</div>
