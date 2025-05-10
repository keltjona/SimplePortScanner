document.addEventListener('DOMContentLoaded', function() {
    const scanButton = document.getElementById('scan-button');
    const targetInput = document.getElementById('target');
    const portsInput = document.getElementById('ports');
    const statusElement = document.getElementById('status');
    const resultsBody = document.getElementById('results-body');
    const summaryElement = document.getElementById('summary');
    const resultsContainer = document.querySelector('.results-container');
    
    scanButton.addEventListener('click', function() {
        // Get input values
        const target = targetInput.value.trim();
        const ports = portsInput.value.trim();
        
        // Validate input
        if (!target) {
            alert('Please enter a target IP address or domain name.');
            return;
        }
        
        if (!ports) {
            alert('Please enter ports to scan.');
            return;
        }
        
        // Clear previous results
        resultsBody.innerHTML = '';
        summaryElement.textContent = '';
        
        // Show results container
        resultsContainer.style.display = 'block';
        
        // Update status
        statusElement.textContent = 'Scanning... Please wait.';
        statusElement.className = 'status scanning';
        
        // Disable scan button during scan
        scanButton.disabled = true;
        
        // Add detailed console logging
        console.log(`Starting scan of target: ${target}, ports: ${ports}`);
        
        // Send scan request
        fetch('/api/scan', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                target: target,
                ports: ports
            })
        })
        .then(response => {
            console.log(`Received response with status: ${response.status}`);
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`Server responded with status ${response.status}: ${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            // Update status
            statusElement.textContent = 'Scan completed!';
            statusElement.className = 'status completed';
            
            // Count open ports
            const openPorts = data.results.filter(result => result.State === 'open');
            
            // Update summary
            summaryElement.textContent = `Scanned ${data.results.length} ports, found ${openPorts.length} open port(s).`;
            
            // Display results
            openPorts.forEach(result => {
                const row = document.createElement('tr');
                
                const portCell = document.createElement('td');
                portCell.textContent = result.Port;
                portCell.className = 'port-open';
                
                const stateCell = document.createElement('td');
                stateCell.textContent = result.State;
                
                const serviceCell = document.createElement('td');
                serviceCell.textContent = result.Service || 'Unknown';
                
                row.appendChild(portCell);
                row.appendChild(stateCell);
                row.appendChild(serviceCell);
                
                resultsBody.appendChild(row);
            });
            
            // Re-enable scan button
            scanButton.disabled = false;
        })
        .catch(error => {
            console.error('Error details:', error);
            statusElement.textContent = `Error: ${error.message}`;
            statusElement.className = 'status error';
            scanButton.disabled = false;
        });
    });
});
