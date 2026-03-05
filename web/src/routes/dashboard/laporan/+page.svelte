<script>
    // Mock data for reports
    const reports = [
        { id: 1, title: 'Laporan Situs Aktif Bulan Mei 2023', date: '2023-05-31', type: 'Statistik', author: 'Ahmad Fauzi', status: 'Selesai' },
        { id: 2, title: 'Laporan Monitoring Situs Mingguan', date: '2023-05-28', type: 'Monitoring', author: 'Siti Nurhaliza', status: 'Selesai' },
        { id: 3, title: 'Laporan Kerusakan Situs Agama', date: '2023-05-25', type: 'Inspeksi', author: 'Budi Santoso', status: 'Selesai' },
        { id: 4, title: 'Laporan Verifikasi Situs Baru', date: '2023-05-20', type: 'Verifikasi', author: 'Lina Marlina', status: 'Selesai' },
        { id: 5, title: 'Laporan Rekapitulasi Bulan April 2023', date: '2023-04-30', type: 'Statistik', author: 'Rizki Pratama', status: 'Selesai' },
        { id: 6, title: 'Laporan Evaluasi Program Kerja', date: '2023-04-15', type: 'Evaluasi', author: 'Ahmad Fauzi', status: 'Selesai' }
    ];

    // Mock data for report statistics
    const reportStats = [
        { title: 'Total Laporan', value: '24', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', color: 'primary' },
        { title: 'Laporan Bulan Ini', value: '6', icon: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z', color: 'secondary' },
        { title: 'Laporan Belum Diproses', value: '0', icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z', color: 'accent' },
        { title: 'Laporan Selesai', value: '24', icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z', color: 'success' }
    ];
</script>

<div class="max-w-7xl mx-auto">
    <!-- Page Header -->
    <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">Laporan dan Dokumentasi</h1>
        <p class="text-gray-600">Kelola dan akses laporan-laporan terkait situs keagamaan</p>
    </div>

    <!-- Report Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {#each reportStats as stat (stat.title)}
            <div class="card bg-base-100 shadow-xl border border-base-200">
                <div class="card-body">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">{stat.title}</p>
                            <h3 class="text-3xl font-bold text-{stat.color}">{stat.value}</h3>
                        </div>
                        <div class={`bg-${stat.color}-100 p-3 rounded-full`}>
                            <svg xmlns="http://www.w3.org/2000/svg" class={`h-6 w-6 text-${stat.color}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{stat.icon}" />
                            </svg>
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <!-- Controls -->
    <div class="card bg-base-100 shadow-xl border border-base-200 mb-6">
        <div class="card-body">
            <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div class="form-control">
                    <input
                        type="text"
                        placeholder="Cari laporan..."
                        class="input input-bordered w-full md:w-80"
                    />
                </div>

                <div class="flex flex-wrap gap-2">
                    <select class="select select-bordered">
                        <option value="all">Semua Jenis</option>
                        <option value="statistik">Statistik</option>
                        <option value="monitoring">Monitoring</option>
                        <option value="inspeksi">Inspeksi</option>
                        <option value="verifikasi">Verifikasi</option>
                        <option value="evaluasi">Evaluasi</option>
                    </select>

                    <button class="btn btn-primary">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        Buat Laporan Baru
                    </button>
                </div>
            </div>
        </div>
    </div>

    <!-- Reports Table -->
    <div class="card bg-base-100 shadow-xl border border-base-200">
        <div class="card-body overflow-x-auto">
            <table class="table table-zebra">
                <thead>
                    <tr>
                        <th>Judul Laporan</th>
                        <th>Tanggal</th>
                        <th>Jenis</th>
                        <th>Dibuat Oleh</th>
                        <th>Status</th>
                        <th>Aksi</th>
                    </tr>
                </thead>
                <tbody>
                    {#each reports as report (report.id)}
                        <tr>
                            <td>
                                <div class="font-medium">{report.title}</div>
                            </td>
                            <td>{report.date}</td>
                            <td>
                                <span class="badge badge-outline">{report.type}</span>
                            </td>
                            <td>{report.author}</td>
                            <td>
                                <span class="badge badge-success">{report.status}</span>
                            </td>
                            <td>
                                <div class="flex space-x-2">
                                    <button class="btn btn-xs btn-outline">Lihat</button>
                                    <button class="btn btn-xs btn-outline btn-primary">Unduh</button>
                                </div>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
</div>
