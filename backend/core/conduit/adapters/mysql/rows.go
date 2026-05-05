package mysql

type result struct {
	rows int64
}

func (r *result) RowsAffected() int64 { return r.rows }
