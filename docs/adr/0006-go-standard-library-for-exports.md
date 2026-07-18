# ADR-0006: Go Standard Library for CSV/PDF Export

**Date**: 2026-07-18  
**Status**: Proposed  
**Drivers**: @development-team

## Context

The clinic management system requires export functionality for:
- **CSV exports** — masked patient data for financial reports
- **PDF receipts** — printable receipt documents
- **Excel exports** — (future requirement for advanced reporting)

The team is evaluating whether Go's standard library is sufficient or if Python's ecosystem (pandas, openpyxl, reportlab) is needed.

## Considered Options

### Option 1: Go Standard Library + Third-Party Packages

**CSV:**
- `encoding/csv` — built-in, robust, handles edge cases

**PDF:**
- `github.com/jung-kurt/gofpdf` — mature, simple API, 1.8M downloads
- `github.com/signintech/gopdf` — active development, good Unicode support

**Excel:**
- `github.com/xuri/excelize/v2` — full Excel feature support, 4.5K stars
- `github.com/tealeg/xlsx` — simpler API, read/write support

**Pros:**
- Single binary deployment maintained
- No Python runtime dependency
- Type-safe compile-time checks
- Performance: Go is 10-100x faster than Python for data processing

**Cons:**
- Less feature-rich than Python ecosystem
- PDF generation more verbose than reportlab
- Team less familiar with Go than Python

### Option 2: Python (Flask/FastAPI) + Python Ecosystem

**CSV:**
- `csv` module — built-in, simple
- `pandas` — powerful data manipulation, Excel export built-in

**PDF:**
- `reportlab` — industry standard, extensive features
- `weasyprint` — HTML/CSS to PDF (easier styling)

**Excel:**
- `openpyxl` — full Excel support
- `pandas` — DataFrame → Excel in 2 lines

**Pros:**
- Rich ecosystem, battle-tested libraries
- Team familiarity with Python
- Rapid prototyping for complex reports
- Better documentation for reporting libraries

**Cons:**
- Python runtime required (25-50MB + venv)
- Dependency management (pip, requirements.txt)
- Slower execution (not critical for exports, but noticeable)
- Deployment complexity (virtualenv, pip install, path management)

### Option 3: Hybrid Approach

- Go backend for core CRUD operations
- Python subprocess for complex report generation
- Shared SQLite database

**Pros:**
- Best of both worlds
- Go for performance-critical paths
- Python for reporting flexibility

**Cons:**
- Architecture complexity
- Two runtimes to manage
- Inter-process communication overhead
- Debugging complexity

## Decision

**Adopt Option 1: Go Standard Library + Third-Party Packages**

### Rationale

1. **Deployment simplicity is paramount** — Single binary is a core requirement for a desktop app targeting non-technical users (clinic practitioners)

2. **Export requirements are modest** — 
   - CSV: Simple tabular data, `encoding/csv` is sufficient
   - PDF: Receipt templates are fixed-layout, gofpdf handles this well
   - Excel: Future requirement, excelize provides full support

3. **Performance matters for large datasets** — If the clinic has 10,000+ receipts, Go's performance will be noticeable for report generation

4. **Team can learn Go** — The team's Python familiarity is a short-term advantage; Go is easier to master for this use case (no async, no decorators, simple concurrency)

5. **No Python-specific features needed** — The exports don't require pandas' data analysis or reportlab's advanced layout features

## Implementation Strategy

### CSV Export (Priority: High)

```go
import "encoding/csv"

func exportCSV(w http.ResponseWriter, receipts []Receipt) error {
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment;filename=reports.csv")
    
    writer := csv.NewWriter(w)
    defer writer.Flush()
    
    // Write header
    writer.Write([]string{"Receipt No", "Patient", "Date", "Amount"})
    
    // Write data (masked)
    for _, r := range receipts {
        writer.Write([]string{
            r.ReceiptNumber,
            maskName(r.PatientName),
            r.VisitDate.Format("2006-01-02"),
            fmt.Sprintf("$%.2f", float64(r.GrandTotal)/100),
        })
    }
    
    return nil
}
```

### PDF Receipt (Priority: High)

```go
import "github.com/jung-kurt/gofpdf"

func generateReceiptPDF(receipt Receipt) ([]byte, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    
    // Clinic header
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(190, 10, clinic.Name)
    
    // Receipt details
    pdf.SetFont("Arial", "", 12)
    pdf.CellFormat(190, 10, fmt.Sprintf("Receipt No: %s", receipt.ReceiptNumber), "", 1, "L", false, 0, "")
    
    // Line items table
    // ... (gofpdf supports tables via CellFormat)
    
    var buf bytes.Buffer
    pdf.Output(&buf)
    return buf.Bytes(), nil
}
```

**Unicode/Chinese Support:**
- Use `AddFontFromFamily()` with custom TTF fonts
- Example: `pdf.AddFontFromFamily("SimSun", "", "simsun.ttf", "")`

### Excel Export (Priority: Low, Future)

```go
import "github.com/xuri/excelize/v2"

func exportExcel(receipts []Receipt) ([]byte, error) {
    f := excelize.NewFile()
    
    // Set headers
    f.SetCellValue("Sheet1", "A1", "Receipt No")
    f.SetCellValue("Sheet1", "B1", "Patient")
    // ...
    
    // Write data
    for i, r := range receipts {
        row := i + 2
        f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), r.ReceiptNumber)
        // ...
    }
    
    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
```

## Consequences

### Positive

- **Single binary maintained** — End users download one `.exe` file, double-click to run
- **No Python dependency hell** — No venv, pip, requirements.txt for end users
- **Faster execution** — Report generation is near-instant vs seconds for Python
- **Type safety** — Compile-time checks prevent export format errors

### Negative

- **Learning curve** — Team needs to learn Go's PDF/Excel libraries
- **Less documentation** — Go libraries have fewer tutorials than Python equivalents
- **More verbose** — PDF layout code is more explicit than reportlab's DSL

### Migration Path

If the team later finds Go's PDF/Excel libraries insufficient:

1. **Prototype in Python first** — Use reportlab/openpyxl to design complex reports
2. **Translate to Go** — Implement the same layout in gofpdf/excelize
3. **Hybrid as last resort** — Only use Python subprocess if a feature is impossible in Go

## When to Revisit

This decision should be revisited if:
- PDF requirements exceed gofpdf capabilities (complex multi-column layouts, advanced typography)
- Excel pivot tables or advanced analytics are required
- Team productivity in Go is significantly lower than expected after 2 weeks of development

## References

- **gofpdf**: https://github.com/jung-kurt/gofpdf (1.8M downloads, 7K stars)
- **excelize**: https://github.com/xuri/excelize (4.5K stars, active maintenance)
- **Go CSV**: https://pkg.go.dev/encoding/csv (standard library)
- **Python reportlab**: https://www.reportlab.com/ (for comparison)

---

**Last updated**: 2026-07-18  
**Maintained by**: Development Team