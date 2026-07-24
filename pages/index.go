package pages

var Index = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Database Entries Dashboard</title>
    <style>
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }
    
    body {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        min-height: 100vh;
        padding: 40px 20px;
    }
    
    .container {
        max-width: 1200px;
        margin: 0 auto;
    }
    
    header {
        text-align: center;
        margin-bottom: 40px;
        color: white;
    }
    
    header h1 {
        font-size: 2.5rem;
        margin-bottom: 10px;
        text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
    }
    
    header p {
        font-size: 1.1rem;
        opacity: 0.9;
    }
    
    .stats-bar {
        display: flex;
        justify-content: space-around;
        flex-wrap: wrap;
        gap: 20px;
        margin-bottom: 30px;
    }
    
    .stat-card {
        background: white;
        padding: 20px 30px;
        border-radius: 12px;
        box-shadow: 0 4px 15px rgba(0,0,0,0.1);
        text-align: center;
        min-width: 150px;
        transition: transform 0.3s ease;
    }
    
    .stat-card:hover {
        transform: translateY(-5px);
    }
    
    .stat-value {
        font-size: 2rem;
        font-weight: bold;
        color: #667eea;
        margin-bottom: 5px;
    }
    
    .stat-label {
        color: #666;
        font-size: 0.9rem;
    }
    
    .entries-table {
        background: white;
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 8px 30px rgba(0,0,0,0.15);
    }
    
    table {
        width: 100%;
        border-collapse: collapse;
    }
    
    thead {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
    }
    
    th {
        padding: 18px 15px;
        text-align: left;
        font-weight: 600;
        font-size: 0.95rem;
    }
    
    td {
        padding: 15px;
        border-bottom: 1px solid #eee;
        color: #333;
    }
    
    tbody tr {
        transition: background-color 0.2s ease;
    }
    
    tbody tr:hover {
        background-color: #f5f7ff;
    }
    
    tbody tr:last-child td {
        border-bottom: none;
    }
    
    .status {
        padding: 5px 12px;
        border-radius: 20px;
        font-size: 0.85rem;
        font-weight: 500;
        display: inline-block;
    }
    
    .status-active {
        background: #e8f5e9;
        color: #2e7d32;
    }
    
    .status-inactive {
        background: #ffebee;
        color: #c62828;
    }
    
    .status-pending {
        background: #fff3e0;
        color: #ef6c00;
    }
    
    .btn {
        background: #667eea;
        color: white;
        border: none;
        padding: 8px 16px;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.9rem;
        transition: background 0.3s ease;
    }
    
    .btn:hover {
        background: #5a6fd6;
    }
    
    .footer {
        text-align: center;
        margin-top: 40px;
        color: white;
        opacity: 0.8;
        font-size: 0.9rem;
    }
    
    @media (max-width: 768px) {
        header h1 {
            font-size: 1.8rem;
        }
    
        .stats-bar {
            flex-direction: column;
        }
    
        table {
            font-size: 0.85rem;
        }
    
        th, td {
            padding: 10px 8px;
        }
    }

    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 Database Entries Dashboard</h1>
            <p>View and manage all stored records</p>
        </header>

        <section class="stats-bar">
            <div class="stat-card">
                <div class="stat-value"></div>
                <div class="stat-label">Total Entries</div>
            </div>
            <div class="stat-card">
                <div class="stat-value"></div>
                <div class="stat-label">Active</div>
            </div>
            <div class="stat-card">
                <div class="stat-value"></div>
                <div class="stat-label">Growth This Month</div>
            </div>
        </section>

        <section class="entries-table">
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Email</th>
                        <th>Date Added</th>
                        <th>Status</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    <!-- Populate with dynamic data -->
                </tbody>
            </table>
        </section>

        <footer class="footer">
            <p>&copy; 2024 Your Application</p>
        </footer>
    </div>
</body>
</html>
`
