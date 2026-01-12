/**
 * Tesouro Direto Google Sheets Custom Function
 * 
 * Fetches Tesouro Direto bond data from the GitHub Pages API
 * 
 * @param {string} nome - Bond name (e.g., "Tesouro IPCA+ 2035")
 * @param {string} campo - Field to return: "data_vencimento", "data_base", "data_inicio", 
 *                         "taxa_compra_manha", "taxa_venda_manha", "pu_compra_manha", 
 *                         "pu_venda_manha", "pu_base_manha"
 * @param {string} data_vencimento - Optional: Maturity date in ISO format (yyyy-mm-dd) 
 *                                   to differentiate bonds with the same name
 * @return {string|number} The requested field value
 * @customfunction
 */
function TESOURODIRETO(nome, campo, data_vencimento) {
  // Validate required parameters
  if (!nome || nome === '') {
    throw new Error('TESOURODIRETO: Nome is required');
  }
  if (!campo || campo === '') {
    throw new Error('TESOURODIRETO: Campo is required');
  }
  
  // Valid fields
  const validFields = [
    'data_vencimento', 'data_base', 'data_inicio', 'data_conversao',
    'taxa_compra_manha', 'taxa_venda_manha',
    'pu_compra_manha', 'pu_venda_manha', 'pu_base_manha'
  ];
  
  if (!validFields.includes(campo)) {
    throw new Error('TESOURODIRETO: Invalid campo. Valid options: ' + validFields.join(', '));
  }
  
  // Replace <username> and <repo> with your actual GitHub Pages URL
  // Example: https://brunompagani.github.io/tesouro-api/latest.json
  const apiUrl = 'https://YOUR_USERNAME.github.io/YOUR_REPO_NAME/latest.json';
  
  try {
    // Fetch JSON data
    const response = UrlFetchApp.fetch(apiUrl);
    const data = JSON.parse(response.getContentText());
    
    // Filter bonds by name
    let matches = data.filter(function(bond) {
      return bond.nome === nome;
    });
    
    // If no matches found
    if (matches.length === 0) {
      throw new Error('TESOURODIRETO: No bond found with nome "' + nome + '"');
    }
    
    // If multiple matches and no date provided
    if (matches.length > 1 && (!data_vencimento || data_vencimento === '')) {
      const dates = matches.map(function(bond) {
        return bond.data_vencimento;
      }).join(', ');
      throw new Error('TESOURODIRETO: Multiple bonds found with nome "' + nome + 
                      '". Please provide data_vencimento parameter. Available dates: ' + dates);
    }
    
    // Filter by date if provided
    if (data_vencimento && data_vencimento !== '') {
      matches = matches.filter(function(bond) {
        return bond.data_vencimento === data_vencimento;
      });
      
      if (matches.length === 0) {
        throw new Error('TESOURODIRETO: No bond found with nome "' + nome + 
                        '" and data_vencimento "' + data_vencimento + '"');
      }
    }
    
    // Get the first (and should be only) match
    const bond = matches[0];
    
    // Return the requested field
    if (bond.hasOwnProperty(campo)) {
      return bond[campo];
    } else {
      throw new Error('TESOURODIRETO: Field "' + campo + '" not found in bond data');
    }
    
  } catch (error) {
    if (error.message && error.message.startsWith('TESOURODIRETO:')) {
      throw error;
    }
    // Otherwise, wrap it
    throw new Error('TESOURODIRETO: Error fetching data - ' + error.message);
  }
}
