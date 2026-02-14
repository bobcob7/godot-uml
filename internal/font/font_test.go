package font

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureText(t *testing.T) {
	t.Parallel()
	t.Run("SingleLine", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("Hello", 13, FamilySans)
		require.NoError(t, err)
		assert.Greater(t, size.Width, 0.0)
		assert.Greater(t, size.Height, 0.0)
	})
	t.Run("EmptyString", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("", 13, FamilySans)
		require.NoError(t, err)
		assert.Equal(t, 0.0, size.Width)
		assert.Greater(t, size.Height, 0.0, "empty string should still have line height")
	})
	t.Run("MultiLine", func(t *testing.T) {
		t.Parallel()
		single, err := MeasureText("Hello", 13, FamilySans)
		require.NoError(t, err)
		multi, err := MeasureText("Hello\nWorld", 13, FamilySans)
		require.NoError(t, err)
		assert.Greater(t, multi.Height, single.Height, "two lines should be taller than one")
		assert.InDelta(t, multi.Height, single.Height*2, 1.0, "two lines should be ~2x one line height")
	})
	t.Run("LongerTextIsWider", func(t *testing.T) {
		t.Parallel()
		short, err := MeasureText("Hi", 13, FamilySans)
		require.NoError(t, err)
		long, err := MeasureText("Hello World", 13, FamilySans)
		require.NoError(t, err)
		assert.Greater(t, long.Width, short.Width)
	})
	t.Run("LargerFontIsLarger", func(t *testing.T) {
		t.Parallel()
		small, err := MeasureText("Hello", 10, FamilySans)
		require.NoError(t, err)
		large, err := MeasureText("Hello", 20, FamilySans)
		require.NoError(t, err)
		assert.Greater(t, large.Width, small.Width)
		assert.Greater(t, large.Height, small.Height)
	})
	t.Run("MonoFamily", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("Hello", 13, FamilyMono)
		require.NoError(t, err)
		assert.Greater(t, size.Width, 0.0)
		assert.Greater(t, size.Height, 0.0)
	})
	t.Run("BoldFamily", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("Hello", 13, FamilyBold)
		require.NoError(t, err)
		assert.Greater(t, size.Width, 0.0)
		assert.Greater(t, size.Height, 0.0)
	})
	t.Run("MonoEqualWidthChars", func(t *testing.T) {
		t.Parallel()
		narrow, err := MeasureText("iiiii", 13, FamilyMono)
		require.NoError(t, err)
		wide, err := MeasureText("MMMMM", 13, FamilyMono)
		require.NoError(t, err)
		assert.InDelta(t, narrow.Width, wide.Width, 1.0, "mono font should have equal char widths")
	})
	t.Run("MultiLineMaxWidth", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("Hi\nHello World", 13, FamilySans)
		require.NoError(t, err)
		longLine, err := MeasureText("Hello World", 13, FamilySans)
		require.NoError(t, err)
		assert.InDelta(t, longLine.Width, size.Width, 0.1, "width should match longest line")
	})
	t.Run("KnownDimensions", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("A", 13, FamilySans)
		require.NoError(t, err)
		assert.Greater(t, size.Width, 5.0)
		assert.Less(t, size.Width, 15.0)
		assert.Greater(t, size.Height, 10.0)
		assert.Less(t, size.Height, 25.0)
	})
	t.Run("ThreeLines", func(t *testing.T) {
		t.Parallel()
		single, err := MeasureText("A", 13, FamilySans)
		require.NoError(t, err)
		triple, err := MeasureText("A\nB\nC", 13, FamilySans)
		require.NoError(t, err)
		assert.InDelta(t, triple.Height, single.Height*3, 1.0)
	})
	t.Run("DefaultFamilyIsSans", func(t *testing.T) {
		t.Parallel()
		size, err := MeasureText("Hello", 13, "unknown")
		require.NoError(t, err)
		sansSize, err := MeasureText("Hello", 13, FamilySans)
		require.NoError(t, err)
		assert.InDelta(t, sansSize.Width, size.Width, 0.1)
	})
}
