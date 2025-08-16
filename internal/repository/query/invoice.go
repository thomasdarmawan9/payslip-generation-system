package query

var CountInvoices = `
SELECT COUNT(*) FROM invoices WHERE deleted_at IS NULL
`

// CreateInvoice inserts invoice only (items inserted separately)
var CreateInvoice = `
INSERT INTO invoices (
	invoice_id,
	subject,
	issue_date,
	due_date,
	status,
	customer_id,
	total_items,
	sub_total,
	tax_percent,
	tax_amount,
	grand_total,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
`

// CreateInvoiceItem inserts one invoice item
var CreateInvoiceItem = `
INSERT INTO invoice_items (
	invoice_id,
	item_name,
	qty,
	unit_price,
	amount,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, ?, NOW(), NOW())
`

// GetInvoiceByID gets one invoice with basic fields (you may join for customer if needed)
var GetInvoiceByID = `
SELECT * FROM invoices WHERE id = ? AND deleted_at IS NULL
`

// GetInvoiceItemsByInvoiceID retrieves all items for a given invoice
var GetInvoiceItemsByInvoiceID = `
SELECT * FROM invoice_items WHERE invoice_id = ? AND deleted_at IS NULL
`

var GetInvoiceByInvoiceID = `
SELECT * FROM invoices WHERE invoice_id = ? AND deleted_at IS NULL
`

// GetAllInvoices retrieves all invoices
var GetAllInvoices = `
SELECT * FROM invoices WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?
`

// UpdateInvoice updates invoice by id
var UpdateInvoice = `
UPDATE invoices SET
	subject = ?,
	issue_date = ?,
	due_date = ?,
	status = ?,
	customer_id = ?,
	total_items = ?,
	sub_total = ?,
	tax_percent = ?,
	tax_amount = ?,
	grand_total = ?,
	updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL
`

// DeleteInvoice soft-deletes invoice by setting deleted_at
var DeleteInvoice = `
UPDATE invoices SET deleted_at = NOW() WHERE id = ?
`

// DeleteInvoiceItems deletes invoice items related to invoice
var DeleteInvoiceItems = `
UPDATE invoice_items SET deleted_at = NOW() WHERE invoice_id = ?
`
